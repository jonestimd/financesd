import agent from 'superagent';
import {AccountModel, ICompany, IAccount} from '../model/AccountModel';
import {indexById, sortByName} from '../model/entityUtils';
import {IMessageStore} from './MessageStore';
import {computed, flow, observable} from 'mobx';

const graphql = `{
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

export class AccountStore {
    private loading: boolean = false;
    @observable
    private companiesById: {[id: string]: ICompany} = {};
    @observable
    private accountsById: {[id: string]: AccountModel} = {};
    private messageStore: IMessageStore;

    constructor(messageStore: IMessageStore) {
        this.messageStore = messageStore;
    }

    @computed
    get accounts(): AccountModel[] {
        return Object.values(this.accountsById).sort(AccountModel.compare);
    }

    @computed
    get companies(): ICompany[] {
        return sortByName(this.companiesById);
    }

    getAccounts(): void {
        if (!this.loading) {
            this.messageStore.addProgressMessage(loadingAccounts);
            this.loadAccounts();
        }
    }

    private loadAccounts = flow(function* () {
        this.loading = true;
        try {
            const {body: {data}}: IAccountsResponse = yield agent.post('/finances/api/v1/graphql').send({query: graphql});
            this.companiesById = indexById(data.companies);
            this.accountsById = indexById(data.accounts.map((account) => new AccountModel(account, this.companiesById[account.companyId])));
        } catch (err) {
            console.error('error gettting accounts', err);
        } finally {
            this.loading = false;
            this.messageStore.removeProgressMessage(loadingAccounts);
        }
    });
}