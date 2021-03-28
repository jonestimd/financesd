import type settingsStore from './settingsStore';

// importing from settingsStore breaks isolateModules
const hideClosedAccounts = 'finances:hideClosedAccounts';
const transactionsView = 'finances:transactionsView';

type SettingStoreModule = {default: typeof settingsStore};

describe('settingsStore', () => {
    let store: typeof settingsStore;
    const loadStore = () => {
        jest.isolateModules(() => {
            store = jest.requireActual<SettingStoreModule>('./settingsStore').default;
        });
    };
    describe('hideClosedAccounts', () => {
        describe('get', () => {
            it('returns false if value does not exist in local storage', () => {
                window.localStorage.removeItem(hideClosedAccounts);

                loadStore();

                expect(store.hideClosedAccounts).toBe(false);
            });
            it('returns true if value exists in local storage', () => {
                window.localStorage.setItem(hideClosedAccounts, 'true');

                loadStore();

                expect(store.hideClosedAccounts).toBe(true);
            });
        });
        describe('set', () => {
            it('adds value to local storage when set to true', () => {
                loadStore();

                store.hideClosedAccounts = true;

                expect(window.localStorage.getItem(hideClosedAccounts)).toEqual('true');
            });
            it('removes value from local storage when set to false', () => {
                window.localStorage.setItem(hideClosedAccounts, 'anything');
                loadStore();

                store.hideClosedAccounts = false;

                expect(window.localStorage.getItem(hideClosedAccounts)).toEqual(null);
            });
        });
    });
    describe('transactionsView', () => {
        describe('get', () => {
            it('returns list if value does not exist in local storage', () => {
                window.localStorage.removeItem(transactionsView);

                loadStore();

                expect(store.transactionsView).toEqual('list');
            });
            it('returns value in local storage', () => {
                window.localStorage.setItem(transactionsView, 'table');

                loadStore();

                expect(store.transactionsView).toEqual('table');
            });
        });
        describe('set', () => {
            it('adds value to local storage when set to table', () => {
                loadStore();

                store.transactionsView = 'table';

                expect(window.localStorage.getItem(transactionsView)).toEqual('table');
            });
            it('removes value from local storage when set to list', () => {
                window.localStorage.setItem(transactionsView, 'table');
                loadStore();

                store.transactionsView = 'list';

                expect(window.localStorage.getItem(transactionsView)).toEqual(null);
            });
        });
    });
});
