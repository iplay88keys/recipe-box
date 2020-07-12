import { AppBar, Button, Container, CssBaseline, Link, Toolbar, Typography } from "@material-ui/core";
import makeStyles from "@material-ui/core/styles/makeStyles";
import React from "react";

const useStyles = makeStyles((theme) => ({
    "@global": {
        ul: {
            margin: 0,
            padding: 0,
            listStyle: "none"
        }
    },
    appBar: {
        borderBottom: `1px solid ${theme.palette.divider}`
    },
    toolbar: {
        flexWrap: "wrap"
    },
    toolbarTitle: {
        flexGrow: 1
    },
    link: {
        margin: theme.spacing(1, 1.5)
    }
}));

export const Navigation = () => {
    const classes = useStyles();

    return (
        <Container component="main" maxWidth={false} disableGutters={true}>
            <CssBaseline/>
            <AppBar position="static" color="default" elevation={0} className={classes.appBar}>
                <Toolbar className={classes.toolbar}>
                    <Typography variant="h6" color="inherit" noWrap className={classes.toolbarTitle}>
                        Recipe Box
                    </Typography>
                    <nav>
                        <Link variant="button" color="textPrimary" href="/" className={classes.link}>
                            Home
                        </Link>
                        <Link variant="button" color="textPrimary" href="/recipes" className={classes.link}>
                            Recipes
                        </Link>
                    </nav>
                    <Button href="/register" color="primary" variant="outlined" className={classes.link}>
                        Register
                    </Button>
                </Toolbar>
            </AppBar>
        </Container>
    );
};
