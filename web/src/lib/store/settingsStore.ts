import {makeObservable, observable, runInAction} from 'mobx';

const hideClosedAccounts = 'finances:hideClosedAccounts';
const transactionsView = 'finances:transactionsView';

export type ViewMode = 'list' | 'table';

class SettingsStore {
    @observable _hideClosedAccounts = !!window.localStorage.getItem(hideClosedAccounts);
    @observable _transactionsView: ViewMode = window.localStorage.getItem(transactionsView) as ViewMode ?? 'list';

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

    get transactionsView() {
        return this._transactionsView;
    }

    set transactionsView(mode: ViewMode) {
        runInAction(() => this._transactionsView = mode);
        if (mode !== 'list') window.localStorage.setItem(transactionsView, mode);
        else window.localStorage.removeItem(transactionsView);
    }
}

export default new SettingsStore();
