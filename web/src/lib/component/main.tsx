import {createBrowserHistory} from 'history';
import React from 'react';
import ReactDom from 'react-dom';
import {Route, Router, Switch} from 'react-router';
import {RootStore, RootStoreContext} from '../store/RootStore';
import ProgressMessage from './ProgressMessage';
import AccountsPage from './AccountsPage';
import TransactionsPage from './transaction/TransactionsPage';
import CategoriesPage from './CategoriesPage';
import {createMuiTheme, ThemeProvider} from '@material-ui/core/styles';

const history = createBrowserHistory();

const NotFound = () => <h4>No such page</h4>;

const myTheme = createMuiTheme({
    palette: {
        primary: {main: '#6200ee'},
        secondary: {main: 'rgb(84, 110, 122)'},
    },
    overrides: {
        MuiTable: {
            root: {
                borderCollapse: 'separate',
            },
        },
        MuiTableCell: {
            root: {
                padding: '1px 3px',
                lineHeight: 1.5,
            },
            head: {
                lineHeight: 1.5,
            },
        },
    },
    props: {
        MuiTypography: {
            variantMapping: {
                body1: 'div',
                body2: 'div',
            },
        },
    },
});

const Routes: React.FC = () => {
    const rootStore = (window.rootStore = window.rootStore || new RootStore());
    React.useEffect(() => {
        rootStore.accountStore.loadAccounts();
        rootStore.categoryStore.loadCategories();
        rootStore.groupStore.loadGroups();
        rootStore.payeeStore.loadPayees();
        rootStore.securityStore.loadSecurities();
    }, []);
    return (
        <RootStoreContext.Provider value={rootStore}>
            <ThemeProvider theme={myTheme}>
                <main className='app-main'>
                    <ProgressMessage />
                    <Router history={history}>
                        <Switch>
                            <Route exact path='/finances' component={AccountsPage} />
                            <Route exact path='/finances/account/:accountId' component={TransactionsPage} />
                            <Route exact path='/finances/categories' component={CategoriesPage} />
                            <Route component={NotFound} />
                        </Switch>
                    </Router>
                </main>
            </ThemeProvider>
        </RootStoreContext.Provider>
    );
};

const domContainer = document.querySelector('#root-container');
ReactDom.render(<Routes />, domContainer);
