import {createBrowserHistory} from 'history';
import React from 'react';
import ReactDom from 'react-dom';
import {Route, Router, Switch} from 'react-router';
import {RootStore, RootStoreContext} from '../store/RootStore';
import ProgressMessage from './ProgressMessage';
import AccountsPage from './AccountsPage';
import TransactionsPage from './transaction/TransactionsPage';
import CategoriesPage from './CategoriesPage';
import PayeesPage from './PayeesPage';
import SecuritiesPage from './SecuritiesPage';
import GroupsPage from './GroupsPage';
import {createMuiTheme, StylesProvider, ThemeProvider} from '@material-ui/core/styles';
import AlertContainer from './AlertContainer';

const history = createBrowserHistory();

const NotFound = () => <h4>No such page</h4>;

const myTheme = createMuiTheme({
    palette: {
        primary: {main: '#6200ee'},
        secondary: {main: 'rgb(84, 110, 122)'},
    },
    zIndex: {
        modal: undefined,
    },
    overrides: {
        MuiTable: {
            root: {
                borderCollapse: 'separate',
            },
        },
        MuiTableCell: {
            root: {
                padding: '2px 3px 1px',
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

export const Routes: React.FC = () => {
    const rootStore = (window.rootStore = window.rootStore || new RootStore());
    React.useEffect(() => {
        void rootStore.accountStore.loadAccounts();
        void rootStore.categoryStore.loadCategories();
        void rootStore.groupStore.loadGroups();
        void rootStore.payeeStore.loadPayees();
        void rootStore.securityStore.loadSecurities();
    }, [rootStore.accountStore, rootStore.categoryStore, rootStore.groupStore, rootStore.payeeStore, rootStore.securityStore]);
    return (
        <RootStoreContext.Provider value={rootStore}>
            <StylesProvider injectFirst>
                <ThemeProvider theme={myTheme}>
                    <main className='app-main'>
                        <ProgressMessage />
                        <AlertContainer />
                        <Router history={history}>
                            <Switch>
                                <Route exact path='/finances' component={AccountsPage} />
                                <Route exact path='/finances/account/:accountId' component={TransactionsPage} />
                                <Route exact path='/finances/categories' component={CategoriesPage} />
                                <Route exact path='/finances/payees' component={PayeesPage} />
                                <Route exact path='/finances/securities' component={SecuritiesPage} />
                                <Route exact path='/finances/groups' component={GroupsPage} />
                                <Route component={NotFound} />
                            </Switch>
                        </Router>
                    </main>
                </ThemeProvider>
            </StylesProvider>
        </RootStoreContext.Provider>
    );
};

const domContainer = document.querySelector('#root-container');
ReactDom.render(<Routes />, domContainer);
