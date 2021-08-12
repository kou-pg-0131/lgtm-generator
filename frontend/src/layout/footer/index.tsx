import React from 'react';
import Link from 'next/link';
import { Routes } from '~/routes';
import ExternalLink from '~/components/externalLink';
import {
  Box,
  Typography,
} from '@material-ui/core';
import {
  createStyles,
  makeStyles,
  Theme,
} from '@material-ui/core/styles';

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
      alignItems: 'center',
      display: 'flex',
      flexDirection: 'column',
      padding: theme.spacing(2),
    },
    listItem: {
      marginBottom: theme.spacing(1),
      textAlign: 'center',
    },
  }),
);

const Footer: React.VFC = () => {
  const classes = useStyles();

  return (
    <Box
      className={classes.root}
      component='footer'
    >
      <Typography component='small'>&copy;2021 koki sato</Typography>

      <ul>
        <li className={classes.listItem}>
          <ExternalLink href='https://github.com/kou-pg-0131/lgtm-generator'>
            View on GitHub
          </ExternalLink>
        </li>
        <li className={classes.listItem}>
          <Link href={Routes.privacyPolicy}>
            <a>
              プライバシーポリシー
            </a>
          </Link>
        </li>
      </ul>
    </Box>
  );
};

export default Footer;
