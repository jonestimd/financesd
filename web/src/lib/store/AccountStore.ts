import * as agent from '../agent';
import {AccountModel, IAccount} from '../model/AccountModel';
import {CompanyModel, ICompany} from '../model/CompanyModel';
import {addToMap, sortValues, sortValuesByName} from '../model/entityUtils';
import {IMessageStore} from './MessageStore';
import {computed, flow, makeObservable, ObservableMap} from 'mobx';
import {LoadResult} from './interfaces';

export const query = `{
    accounts {
        id name type accountNo description closed companyId version transactionCount balance
    }
    companies {
        id name version
    }
}`;

type AccountsResponse = agent.IGraphqlResponse<{accounts: IAccount[], companies: ICompany[]}>;

export const loadingAccounts = 'Loading accounts...';

export default class AccountStore {
    private loading = false;
    private companiesById = new ObservableMap<string, CompanyModel>();
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
    get filteredAccounts() {
        return this.accounts.filter((account) => !account.hide);
    }

    @computed
    get accountsWithoutCompany() {
        return this.filteredAccounts.filter((account) => !account.companyId);
    }

    @computed
    get companies(): CompanyModel[] {
        return sortValuesByName(this.companiesById);
    }

    @computed
    get filteredCompanies() {
        return this.companies.filter((company) => company.filteredAccounts.length > 0);
    }

    getAccount(id?: string | number) {
        return this.accountsById.get('' + id);
    }

    loadAccounts(): Promise<void> | undefined {
        if (!this.loading && Object.keys(this.accounts).length === 0) {
            this.messageStore.addProgressMessage(loadingAccounts);
            return this._loadAccounts();
        }
    }

    private _loadAccounts = flow(function* (this: AccountStore): LoadResult<AccountsResponse> {
        this.loading = true;
        try {
            const {data: {accounts, companies}} = yield agent.graphql('/finances/api/v1/graphql', query);
            addToMap(this.companiesById, companies.map((company) => new CompanyModel(company)));
            addToMap(this.accountsById, accounts.map((account) => new AccountModel(account, this.companiesById.get('' + account.companyId))));
            for (const account of this.accounts) {
                account.company?.accounts.push(account);
            }
        } catch (err) {
            console.error('error gettting accounts', err);
        } finally {
            this.loading = false;
            this.messageStore.removeProgressMessage(loadingAccounts);
        }
    });
}
