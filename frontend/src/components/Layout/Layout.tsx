import React, { useMemo } from 'react';
import Head from 'next/head';
import { RecoilRoot } from 'recoil';
import Header from './Header';
import Footer from './Footer';
import ToastProvider from '~/components/providers/ToastProvider';
import {
  Box,
  Container,
  CssBaseline,
  ThemeProvider,
  StyledEngineProvider,
} from '@mui/material';
import { Theme } from '@mui/material/styles';
import { theme } from './theme';

declare module '@mui/styles/defaultTheme' {
  // eslint-disable-next-line @typescript-eslint/no-empty-interface
  interface DefaultTheme extends Theme {}
}

type LayoutProps = {
  children: React.ReactNode;
  title?: string;
};

const Layout: React.VFC<LayoutProps> = React.memo(props => {
  return (
    <RecoilRoot>
      <StyledEngineProvider injectFirst>
        <ThemeProvider theme={theme}>
          <ToastProvider>
            <LayoutContent {...props} />
          </ToastProvider>
        </ThemeProvider>
      </StyledEngineProvider>
    </RecoilRoot>
  );
});

Layout.displayName = 'Layout';

const LayoutContent: React.VFC<LayoutProps> = React.memo(props => {
  const { title } = props;

  const pageTitle = useMemo(() => {
    const baseTitle = 'LGTM Generator';
    if (!title) {
      return baseTitle;
    }
    return `${title} | ${baseTitle}`;
  }, [title]);

  return (
    <Box
      sx={{
        backgroundColor: theme => theme.palette.primary.light,
        minHeight: '100vh',
      }}
    >
      <Head>
        <title>{pageTitle}</title>
      </Head>
      <CssBaseline />
      <Header />
      <Container component='main' maxWidth='lg' sx={{ pt: 2 }}>
        {props.children}
      </Container>
      <Footer />
    </Box>
  );
});

LayoutContent.displayName = 'LayoutContent';

export default Layout;
