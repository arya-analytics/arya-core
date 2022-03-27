import { createContext, useContext, useEffect } from "react";
import { usePersistedState } from "../Hooks/usePersistedState";
import {
  createTheme as createMaterialTheme,
  PaletteMode,
  Switch,
  SwitchProps,
  ThemeOptions,
  ThemeProvider as MaterialThemeProvider,
} from "@mui/material";
import merge from "lodash.merge";

const defaultTheme = "dark";
const persistedThemeKey = "aryaTheme";

const PaletteContext = createContext<{
  palette: PaletteMode;
  theme: ThemeOptions;
  setPalette: (palette: PaletteMode) => void;
}>({
  palette: "light",
  theme: {},
  setPalette: () => {},
});

export const useThemeContext = () => useContext(PaletteContext);

export const ThemeProvider = ({ children }: React.PropsWithChildren<any>) => {
  const [palette, setPalette] = usePersistedState<PaletteMode>({
    key: persistedThemeKey,
    defaultValue: defaultTheme,
  });
  const theme = getTheme(palette);
  useEffect(() => {
    console.log(theme);
    document.body.style.backgroundColor = theme.palette?.background
      ?.paper as string;
  }, [palette]);
  return (
    <MaterialThemeProvider theme={theme}>
      <PaletteContext.Provider value={{ palette, theme, setPalette }}>
        {children}
      </PaletteContext.Provider>
    </MaterialThemeProvider>
  );
};

const getTheme = (palette: PaletteMode): ThemeOptions => {
  return themes[palette];
};

declare module "@mui/material/styles" {
  interface Theme {
    status: {
      danger: string;
    };
  }

  // allow configuration using `createTheme`
  interface ThemeOptions {
    status?: {
      danger?: string;
    };
  }
}

const baseTheme: ThemeOptions = {
  palette: {
    primary: {
      main: "#3774D0",
    },
  },
  components: {
    MuiTabs: {
      defaultProps: {
        sx: {
          borderBottom: 1,
          borderColor: "divider",
          minHeight: 0,
          "& .MuiTabs-indicator": {
            backgroundColor: "",
          },
          "& .MuiButtonBase-root": {
            height: 36,
            minHeight: 0,
            textTransform: "none",
            color: "text.primary",
          },
        },
      },
    },
    MuiTypography: {
      defaultProps: {
        color: "text.primary",
      },
    },
  },
};

const lightTheme = {
  palette: {
    mode: "light" as PaletteMode,
    secondary: {
      main: "#212121",
    },
    text: {
      primary: "#212121",
    },
  },
};

const darkTheme: ThemeOptions = {
  palette: {
    mode: "dark" as PaletteMode,
    background: {
      default: "#1F1F1F",
    },
    secondary: {
      main: "#e0e0e0",
    },
    text: {
      primary: "#e0e0e0",
    },
  },
};

const createTheme = (theme: ThemeOptions): ThemeOptions => {
  return createMaterialTheme(merge(baseTheme, theme));
};

const themes = {
  light: createTheme(lightTheme),
  dark: createTheme(darkTheme),
};

export const ToggleThemeSwitch = (props: SwitchProps) => {
  const { palette, setPalette } = useThemeContext();
  return (
    <Switch
      onChange={() => setPalette(palette == "light" ? "dark" : "light")}
      checked={palette == "light"}
      {...props}
    />
  );
};
