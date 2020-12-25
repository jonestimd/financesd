import * as agent from '../agent';
import {AccountModel, ICompany, IAccount} from '../model/AccountModel';
import {addToMap, sortValues, sortValuesByName} from '../model/entityUtils';
import {IMessageStore} from './MessageStore';
import {computed, flow, makeObservable, ObservableMap} from 'mobx';
import {LoadResult} from './interfaces';

const query = `{
    accounts {
        id name type accountNo description closed companyId version transactionCount balance
    }
    companies {
        id name version
    }
}`;

type AccountsResponse = agent.IGraphqlResponse<{accounts: IAccount[], companies: ICompany[]}>;

const loadingAccounts = 'Loading accounts...';

export default class AccountStore {
    private loading = false;
    private companiesById = new ObservableMap<string, ICompany>();
    private accountsById = new ObservableMap<string, AccountModel>();
    private messageStore: IMessageStore;

    constructor(messageStore: IMessageStore) {
        makeObservable(this);
        this.messageStore = messageStore;
    }

    @computed
    get accounts(): AccountModel[] {
        return sortValues(this.accountsById, AccountModel.compare);
    }

    @computed
    get companies(): ICompany[] {
        return sortValuesByName(this.companiesById);
    }

    getAccount(id: string | number) {
        return this.accountsById.get('' + id) || {} as AccountModel;
    }

    loadAccounts(): void {
        if (!this.loading && Object.keys(this.accounts).length === 0) {
            this.messageStore.addProgressMessage(loadingAccounts);
            void this._loadAccounts();
        }
    }

    private _loadAccounts = flow(function* (this: AccountStore): LoadResult<AccountsResponse> {
        this.loading = true;
        try {
            const {data} = yield agent.graphql('/finances/api/v1/graphql', query);
            addToMap(this.companiesById, data.companies);
            addToMap(this.accountsById, data.accounts.map((account) => new AccountModel(account, this.companiesById.get('' + account.companyId))));
        } catch (err) {
            console.error('error gettting accounts', err);
        } finally {
            this.loading = false;
            this.messageStore.removeProgressMessage(loadingAccounts);
        }
    });
}
