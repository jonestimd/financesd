import * as agent from '../agent';
import {AccountModel, IAccount} from '../model/account/AccountModel';
import {CompanyModel, ICompany} from '../model/account/CompanyModel';
import {addToMap, sortValues, sortValuesByName} from '../model/entityUtils';
import {IMessageStore} from './MessageStore';
import {computed, flow, makeObservable, ObservableMap} from 'mobx';
import {LoadResult} from './interfaces';
import Loader from './Loader';

const accountFields = `
fragment accountFields on account {
    id name type accountNo description closed companyId version transactionCount balance
}`;

const companyFields = 'fragment companyFields on company {id name version}';

export const query = `${companyFields}
${accountFields}
{
    accounts {...accountFields}
    companies {...companyFields}
}`;

export const updateCompaniesQuery = `${companyFields}
mutation update($add: [String!], $delete: [Int!], $update: [companyInput!]) {
    companies: updateCompanies(add: $add, delete: $delete, update: $update) {...companyFields}
}`;

type AccountsResponse = agent.IGraphqlResponse<{accounts: IAccount[], companies: ICompany[]}>;

export interface IUpdateCompanies {
    add: string[];
    delete: number[];
    update: Pick<ICompany, 'id' | 'name' | 'version'>[];
}

export const loadingAccounts = 'Loading accounts...'; // TODO localization
export const savingCompanies = 'Saving companies...';

export default class AccountStore {
    private loading = false;
    private companiesById = new ObservableMap<string, CompanyModel>();
    private accountsById = new ObservableMap<string, AccountModel>();
    private messageStore: IMessageStore;
    private loader: Loader;

    constructor(messageStore: IMessageStore) {
        makeObservable(this);
        this.messageStore = messageStore;
        this.loader = new Loader(messageStore);
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
            const {data: {accounts, companies}} = yield agent.graphql(query);
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

    saveCompanies(variables: IUpdateCompanies) {
        return this.loader.load<{companies: ICompany[]}>(savingCompanies, {query: updateCompaniesQuery, variables, updater: ({companies}) => {
            variables.delete.forEach((id) => this.companiesById.delete(`${id}`));
            companies.forEach((c) => {
                if (this.companiesById.has(c.id)) this.companiesById.get(c.id)?.update(c);
                else this.companiesById.set(c.id, new CompanyModel(c));
            });
        }});
    }
}
