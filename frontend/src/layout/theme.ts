import { createTheme } from '@material-ui/core';

export const theme = createTheme({
  overrides: {
    MuiCssBaseline: {
      '@global': {
        a: {
          color: 'inherit',
          textDecoration: 'none',
        },
        ul: {
          listStyle: 'none',
          padding: 0,
        },
      },
    },
  },
});
