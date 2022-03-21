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

const PaletteContext = createContext<[PaletteMode, (val: PaletteMode) => void]>(
  [defaultTheme, () => {}]
);

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
      <PaletteContext.Provider value={[palette, setPalette]}>
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
};

const lightTheme = {
  palette: {
    mode: "light" as PaletteMode,
    secondary: {
      main: "#212121",
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
      main: "#e0e0e0"
    }
  },
  components: {
    MuiTypography: {
      defaultProps: {
        color: "text.primary",
      },
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
  const [theme, setTheme] = useThemeContext();
  return (
    <Switch
      onChange={() => setTheme(theme == "light" ? "dark" : "light")}
      checked={theme == "light"}
      {...props}
    />
  );
};
