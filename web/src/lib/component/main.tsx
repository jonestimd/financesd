import {createBrowserHistory} from 'history';
import React from 'react';
import ReactDom from 'react-dom';
import {Route, Router, Switch} from 'react-router';
import {RootStore, RootStoreContext} from '../store/RootStore';
import ProgressMessage from './ProgressMessage';
import AccountsPage from './AccountsPage';

const history = createBrowserHistory();

const NotFound = () => <h4>No such page</h4>;

interface IMainProps {
    match: {url: string};
}

const Main: React.FC<IMainProps> = ({match}) => (
    <main className='app-main'>
        <Switch>
            <Route exact={true} path={match.url} component={AccountsPage}/>
            <Route component={NotFound}/>
        </Switch>
    </main>
);

const Routes: React.FC<{}> = () => {
    const rootStore = (window.rootStore = window.rootStore || new RootStore());
    return (
        <RootStoreContext.Provider value={rootStore}>
            <ProgressMessage/>
            <Router history={history}>
                <Route path='/finances' component={Main}/>
            </Router>
        </RootStoreContext.Provider>
    );
};

const domContainer = document.querySelector('#root-container');
ReactDom.render(<Routes/>, domContainer);