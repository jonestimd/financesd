import {createBrowserHistory} from 'history';
import React from 'react';
import ReactDom from 'react-dom';
import {Route, Router, Switch} from 'react-router';
import {RootStore, RootStoreContext} from '../store/RootStore';
import ProgressMessage from './ProgressMessage';
import AccountsPage from './AccountsPage';
import TransactionsPage from './TransactionsPage';
import CategoriesPage from './CategoriesPage';

const history = createBrowserHistory();

const NotFound = () => <h4>No such page</h4>;

const Routes: React.FC<{}> = () => {
    const rootStore = (window.rootStore = window.rootStore || new RootStore());
    return (
        <RootStoreContext.Provider value={rootStore}>
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
        </RootStoreContext.Provider>
    );
};

const domContainer = document.querySelector('#root-container');
ReactDom.render(<Routes />, domContainer);