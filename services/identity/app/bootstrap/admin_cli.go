package bootstrap

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/mail"
	"os"
	"strconv"
	"strings"
)

const (
	defaultAdminGender   = 1
	defaultAdminLogin    = "admin"
	defaultAdminName     = "Admin"
	defaultAdminLastName = "Admin"
	defaultAdminRole     = "admin"
	defaultAdminPass     = "admin123"
)

func runAdmin(ctx context.Context, args []string, in io.Reader, out io.Writer) error {
	input, err := resolveAdminInput(args, in, out)
	if err != nil {
		return err
	}

	plainPassword := input.Password
	if plainPassword == "" {
		plainPassword = defaultAdminPass
	}
	input.Password = plainPassword

	adminID, created, err := EnsureAdmin(ctx, input)
	if err != nil {
		return err
	}

	action := "updated"
	if created {
		action = "created"
	}

	if _, err := fmt.Fprintf(out, "Admin %s: id=%s email=%s role=%s gender=%d\n", action, adminID, input.Email, input.Role, input.Gender); err != nil {
		return err
	}
	if plainPassword == defaultAdminPass {
		if _, err := fmt.Fprintf(out, "Admin password: %s (override via ADMIN_PASSWORD)\n", plainPassword); err != nil {
			return err
		}
	} else {
		if _, err := fmt.Fprintln(out, "Admin password taken from --password or ADMIN_PASSWORD"); err != nil {
			return err
		}
	}

	return nil
}

func resolveAdminInput(args []string, in io.Reader, out io.Writer) (AdminSetupInput, error) {
	reader := bufio.NewReader(in)

	f := flag.NewFlagSet("admin", flag.ContinueOnError)
	f.SetOutput(io.Discard)

	defaultGender := envOrDefault("ADMIN_GENDER_CODE", strconv.Itoa(defaultAdminGender))

	emailFlag := f.String("email", strings.TrimSpace(os.Getenv("ADMIN_EMAIL")), "")
	loginFlag := f.String("login", envOrDefault("ADMIN_LOGIN", defaultAdminLogin), "")
	nameFlag := f.String("name", envOrDefault("ADMIN_NAME", defaultAdminName), "")
	lastNameFlag := f.String("last-name", envOrDefault("ADMIN_LAST_NAME", defaultAdminLastName), "")
	roleFlag := f.String("role", envOrDefault("ADMIN_ROLE", defaultAdminRole), "")
	genderFlag := f.String("gender-code", defaultGender, "")
	passwordFlag := f.String("password", strings.TrimSpace(os.Getenv("ADMIN_PASSWORD")), "")

	if err := f.Parse(args); err != nil {
		return AdminSetupInput{}, fmt.Errorf("parse admin args: %w", err)
	}

	if positional := f.Args(); len(positional) > 0 && strings.TrimSpace(*emailFlag) == "" {
		*emailFlag = strings.TrimSpace(positional[0])
	}

	email := strings.TrimSpace(*emailFlag)
	if email == "" {
		if _, err := fmt.Fprint(out, "Enter the admin email: "); err != nil {
			return AdminSetupInput{}, err
		}
		line, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return AdminSetupInput{}, fmt.Errorf("read admin email: %w", err)
		}
		email = strings.TrimSpace(line)
	}

	if email == "" {
		return AdminSetupInput{}, errors.New("admin email is required")
	}
	if err := validateEmail(email); err != nil {
		return AdminSetupInput{}, err
	}

	login := strings.TrimSpace(*loginFlag)
	if login == "" {
		return AdminSetupInput{}, errors.New("admin login is required")
	}

	name := strings.TrimSpace(*nameFlag)
	if name == "" {
		return AdminSetupInput{}, errors.New("admin name is required")
	}

	lastName := strings.TrimSpace(*lastNameFlag)
	if lastName == "" {
		return AdminSetupInput{}, errors.New("admin last name is required")
	}

	role := strings.TrimSpace(*roleFlag)
	if role == "" {
		return AdminSetupInput{}, errors.New("admin role is required")
	}

	gender, err := strconv.Atoi(strings.TrimSpace(*genderFlag))
	if err != nil {
		return AdminSetupInput{}, fmt.Errorf("invalid gender code: %w", err)
	}
	if gender != 1 && gender != 2 {
		return AdminSetupInput{}, errors.New("admin gender code must be 1 or 2")
	}

	password := strings.TrimSpace(*passwordFlag)

	return AdminSetupInput{
		Login:    login,
		Email:    email,
		Name:     name,
		LastName: lastName,
		Role:     role,
		Gender:   gender,
		Password: password,
	}, nil
}

func validateEmail(email string) error {
	addr, err := mail.ParseAddress(email)
	if err != nil || !strings.EqualFold(strings.TrimSpace(addr.Address), email) {
		return fmt.Errorf("invalid admin email: %q", email)
	}
	return nil
}

func envOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
