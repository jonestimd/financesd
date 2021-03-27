import {makeObservable, observable, runInAction} from 'mobx';

const hideClosedAccounts = 'hideClosedAccounts';

class SettingsStore {
    @observable _hideClosedAccounts = !!window.localStorage.getItem(hideClosedAccounts);

    constructor() {
        makeObservable(this);
    }

    get hideClosedAccounts() {
        return this._hideClosedAccounts;
    }

    set hideClosedAccounts(hide: boolean) {
        runInAction(() => this._hideClosedAccounts = hide);
        if (hide) window.localStorage.setItem(hideClosedAccounts, 'true');
        else window.localStorage.removeItem(hideClosedAccounts);
    }
}

export default new SettingsStore();
