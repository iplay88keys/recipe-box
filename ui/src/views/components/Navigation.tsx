import { AppBar, Button, Container, CssBaseline, Link, Toolbar } from "@material-ui/core";
import makeStyles from "@material-ui/core/styles/makeStyles";
import React from "react";
import { Link as RouterLink } from "react-router-dom";

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

interface NavigationProps {
    loggedIn: boolean
}

export const Navigation = ({loggedIn}: NavigationProps) => {
    const classes = useStyles();

    return (
        <Container component="main" maxWidth={false} disableGutters={true}>
            <CssBaseline/>
            <AppBar position="static" color="default" elevation={0} className={classes.appBar}>
                <Toolbar className={classes.toolbar}>
                    <div className={classes.toolbarTitle}>
                        <Link to="/" component={RouterLink} color="textPrimary" underline="none" variant="h6">
                            My Recipe Library
                        </Link>
                    </div>
                    {loggedIn &&
                    <nav>
                        <Button to="/recipes" component={RouterLink} color="primary" className={classes.link}>
                            Recipes
                        </Button>
                    </nav>
                    }
                    {!loggedIn &&
                    <div>
                        <Button to="/register" component={RouterLink} color="primary" variant="outlined"
                                className={classes.link}>
                            Register
                        </Button>
                        <Button to="/login" component={RouterLink} color="primary" variant="contained"
                                className={classes.link}>
                            Login
                        </Button>
                    </div>
                    }
                    {loggedIn &&
                    <Button to="/login" component={RouterLink} color="secondary" variant="contained"
                            className={classes.link}>
                        Logout
                    </Button>
                    }
                </Toolbar>
            </AppBar>
        </Container>
    );
};
