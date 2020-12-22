import agent from 'superagent';
import {AccountModel, ICompany, IAccount} from '../model/AccountModel';
import {addToMap, sortValues, sortValuesByName} from '../model/entityUtils';
import {IMessageStore} from './MessageStore';
import {computed, flow, makeObservable, ObservableMap} from 'mobx';

const query = `{
    accounts {
        id name type accountNo description closed companyId version transactionCount balance
    }
    companies {
        id name version
    }
}`;

interface IAccountsResponse {
    body: {data: {accounts: IAccount[], companies: ICompany[]}};
}

const loadingAccounts = 'Loading accounts...';

export default class AccountStore {
    private loading: boolean = false;
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
            this._loadAccounts();
        }
    }

    private _loadAccounts = flow(function* () {
        this.loading = true;
        try {
            const {body: {data}}: IAccountsResponse = yield agent.post('/finances/api/v1/graphql').send({query});
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