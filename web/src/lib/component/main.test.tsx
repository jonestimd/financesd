import React from 'react';
import {shallow} from 'enzyme';
import {mockUseEffect} from 'src/test/mockHooks';
import {RootStore, RootStoreContext} from '../store/RootStore';
import * as history from 'history';
import type * as main from './main';
import ProgressMessage from './ProgressMessage';
import {Route, Router} from 'react-router';
import AccountsPage from './AccountsPage';
import TransactionsPage from './transaction/TransactionsPage';
import CategoriesPage from './CategoriesPage';
import PayeesPage from './PayeesPage';
import SecuritiesPage from './SecuritiesPage';
import GroupsPage from './GroupsPage';

jest.mock('../store/RootStore');

describe('main', () => {
    const {RootStore: RealRootStore} = jest.requireActual<{RootStore: typeof RootStore}>('../store/RootStore');
    const rootStore = new RealRootStore();
    const {accountStore, categoryStore, groupStore, payeeStore, securityStore} = rootStore;
    const browserHistory = {location: {pathName: 'here'}, listen: jest.fn()};

    describe('Routes', () => {
        beforeEach(() => {
            const mockRootStore = RootStore as jest.MockedClass<typeof RootStore>;
            mockRootStore.mockReturnValue(rootStore);
            jest.spyOn(accountStore, 'loadAccounts').mockResolvedValue(true);
            jest.spyOn(categoryStore, 'loadCategories').mockResolvedValue(true);
            jest.spyOn(groupStore, 'loadGroups').mockResolvedValue(true);
            jest.spyOn(payeeStore, 'loadPayees').mockResolvedValue(true);
            jest.spyOn(securityStore, 'loadSecurities').mockResolvedValue(true);
            mockUseEffect();
            const container = document.createElement('div');
            jest.spyOn(document, 'querySelector').mockReturnValue(container);
            jest.spyOn(history, 'createBrowserHistory').mockReturnValue(browserHistory as unknown as history.History<unknown>);
        });
        jest.isolateModules(() => {
            it('creates root store and loads stores', () => {
                const {Routes} = jest.requireActual<typeof main>('./main');

                const component = shallow(<Routes />);

                expect(RootStore).toBeCalledTimes(1);
                expect(component.find(RootStoreContext.Provider)).toHaveProp('value', rootStore);
                expect(window.rootStore).toBe(rootStore);
                expect(accountStore.loadAccounts).toBeCalledTimes(1);
                expect(categoryStore.loadCategories).toBeCalledTimes(1);
                expect(groupStore.loadGroups).toBeCalledTimes(1);
                expect(payeeStore.loadPayees).toBeCalledTimes(1);
                expect(securityStore.loadSecurities).toBeCalledTimes(1);
            });
        });
        jest.isolateModules(() => {
            it('reuses existing root store', () => {
                window.rootStore = rootStore;
                const {Routes} = jest.requireActual<typeof main>('./main');

                const component = shallow(<Routes />);

                expect(RootStore).not.toBeCalled();
                expect(component.find(RootStoreContext.Provider)).toHaveProp('value', rootStore);
                expect(window.rootStore).toBe(rootStore);
            });
        });
        jest.isolateModules(() => {
            it('renders progress indicator and routes', () => {
                window.rootStore = rootStore;
                const {Routes} = jest.requireActual<typeof main>('./main');

                const component = shallow(<Routes />);

                expect(component.find(ProgressMessage)).toExist();
                expect(component.find(Router)).toHaveProp('history', browserHistory);
                const routes = component.find(Route);
                expect(routes).toHaveLength(7);
                expect(routes.at(0)).toHaveProps({exact: true, path: '/finances', component: AccountsPage});
                expect(routes.at(1)).toHaveProps({exact: true, path: '/finances/account/:accountId', component: TransactionsPage});
                expect(routes.at(2)).toHaveProps({exact: true, path: '/finances/categories', component: CategoriesPage});
                expect(routes.at(3)).toHaveProps({exact: true, path: '/finances/payees', component: PayeesPage});
                expect(routes.at(4)).toHaveProps({exact: true, path: '/finances/securities', component: SecuritiesPage});
                expect(routes.at(5)).toHaveProps({exact: true, path: '/finances/groups', component: GroupsPage});
            });
        });
    });
});
