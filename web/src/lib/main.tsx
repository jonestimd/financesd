import React from 'react';
import ReactDom from 'react-dom';
import {Route, Router, Switch} from 'react-router';
import {createBrowserHistory} from 'history';
import {RootStoreContext, RootStore} from './stores/RootStore';

const history = createBrowserHistory();

const Hello = () => <h4>Hello</h4>;

const NotFound = () => <h4>No such page</h4>;

interface IMainProps {
    match: {url: string};
}

const Main: React.FC<IMainProps> = ({match}) => (
    <main className="app-main">
        <Switch>
            <Route exact={true} path={match.url} component={Hello}/>
            {/* <Route exact={true} path={`${match.url}/accounts`} component={AccountsPage}/> */}
            <Route component={NotFound}/>
        </Switch>
    </main>
);

const Routes: React.FC<{}> = () => {
    const rootStore = (window.rootStore = window.rootStore || new RootStore());
    return (
        <RootStoreContext.Provider value={rootStore}>
            <Router history={history}>
                <Route path='/finances' component={Main}/>
            </Router>
        </RootStoreContext.Provider>
    );
};

const domContainer = document.querySelector('#root-container');
ReactDom.render(<Routes/>, domContainer);