import type settingsStore from './settingsStore';

type SettingStoreModule = {default: typeof settingsStore};

describe('settingsStore', () => {
    describe('get hideClosedAccounts', () => {
        it('returns false is value does not exist in local storage', () => {
            jest.isolateModules(() => {
                window.localStorage.removeItem('hideClosedAccounts');

                const store = jest.requireActual<SettingStoreModule>('./settingsStore').default;

                expect(store.hideClosedAccounts).toBe(false);
            });
        });
        it('returns true is value exists in local storage', () => {
            jest.isolateModules(() => {
                window.localStorage.setItem('hideClosedAccounts', 'true');

                const store = jest.requireActual<SettingStoreModule>('./settingsStore').default;

                expect(store.hideClosedAccounts).toBe(true);
            });
        });
    });
    describe('set hideClosedAccounts', () => {
        it('adds value to local storage when set to true', () => {
            jest.isolateModules(() => {
                const store = jest.requireActual<SettingStoreModule>('./settingsStore').default;

                store.hideClosedAccounts = true;

                expect(window.localStorage.getItem('hideClosedAccounts')).toEqual('true');
            });
        });
        it('removes value to local storage when set to false', () => {
            jest.isolateModules(() => {
                window.localStorage.setItem('hideClosedAccounts', 'anything');
                const store = jest.requireActual<SettingStoreModule>('./settingsStore').default;

                store.hideClosedAccounts = false;

                expect(window.localStorage.getItem('hideClosedAccounts')).toEqual(null);
            });
        });
    });
});
