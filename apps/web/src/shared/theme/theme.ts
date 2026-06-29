import { createTheme } from "@mui/material/styles";

export const appTheme = createTheme({
  cssVariables: true,
  palette: {
    mode: "dark",
    primary: {
      main: "#90caf9",
      light: "#e3f2fd",
      dark: "#42a5f5",
      contrastText: "#0a1929",
    },
    secondary: {
      main: "#ce93d8",
      light: "#f3e5f5",
      dark: "#ab47bc",
      contrastText: "#16091b",
    },
    error: {
      main: "#f44336",
    },
    warning: {
      main: "#ffa726",
    },
    info: {
      main: "#29b6f6",
    },
    success: {
      main: "#66bb6a",
    },
    background: {
      default: "#0d1117",
      paper: "#131a24",
    },
    text: {
      primary: "rgba(255, 255, 255, 0.87)",
      secondary: "rgba(255, 255, 255, 0.6)",
    },
    divider: "rgba(255, 255, 255, 0.12)",
  },
  shape: {
    borderRadius: 24,
  },
  typography: {
    fontFamily: '"Manrope", "Segoe UI", "Helvetica Neue", Arial, sans-serif',
    h1: {
      fontWeight: 800,
      letterSpacing: "-0.04em",
    },
    h2: {
      fontWeight: 700,
      letterSpacing: "-0.03em",
    },
    h3: {
      fontWeight: 700,
      letterSpacing: "-0.02em",
    },
    button: {
      fontWeight: 700,
      textTransform: "none",
    },
  },
  components: {
    MuiCssBaseline: {
      styleOverrides: {
        ":root": {
          colorScheme: "dark",
        },
        body: {
          backgroundColor: "#0d1117",
          backgroundImage:
            "radial-gradient(circle at top left, rgba(144, 202, 249, 0.16), transparent 28%), radial-gradient(circle at top right, rgba(102, 187, 106, 0.12), transparent 24%), linear-gradient(180deg, #111827 0%, #0d1117 100%)",
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          boxShadow: "0 20px 60px rgba(0, 0, 0, 0.24)",
          backgroundColor: "#131a24",
          backgroundImage: "none",
          border: "1px solid rgba(255, 255, 255, 0.06)",
        },
      },
    },
    MuiButton: {
      defaultProps: {
        disableElevation: true,
      },
      styleOverrides: {
        root: {
          borderRadius: 999,
          paddingInline: 18,
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          backgroundImage: "none",
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          fontWeight: 700,
        },
      },
    },
    MuiAppBar: {
      styleOverrides: {
        root: {
          backgroundImage: "none",
        },
      },
    },
    MuiDrawer: {
      styleOverrides: {
        paper: {
          backgroundImage: "none",
        },
      },
    },
    MuiTableCell: {
      styleOverrides: {
        head: {
          color: "rgba(255, 255, 255, 0.72)",
          fontWeight: 700,
        },
      },
    },
    MuiOutlinedInput: {
      styleOverrides: {
        root: {
          backgroundColor: "rgba(255, 255, 255, 0.03)",
        },
      },
    },
  },
});
