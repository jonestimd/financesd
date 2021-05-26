import {AccountModel, IAccount} from '../model/account/AccountModel';
import {CompanyModel, ICompany} from '../model/account/CompanyModel';
import {addToMap, sortValues, sortValuesByName} from '../model/entityUtils';
import {IMessageStore} from './MessageStore';
import {computed, makeObservable, ObservableMap} from 'mobx';
import Loader from './Loader';
import AlertStore from './AlertStore';

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
mutation update($add: [String!], $delete: [idVersion!], $update: [companyInput!]) {
    companies: updateCompanies(add: $add, delete: $delete, update: $update) {...companyFields}
}`;

interface IAccountsResponse {
    accounts: IAccount[];
    companies: ICompany[];
}

export interface IUpdateCompanies {
    add: string[];
    delete: Pick<ICompany, 'id' | 'version'>[];
    update: Pick<ICompany, 'id' | 'name' | 'version'>[];
}

export const loadingAccounts = 'Loading accounts';
export const savingCompanies = 'Saving companies';

export default class AccountStore {
    private loading = false;
    private companiesById = new ObservableMap<number, CompanyModel>();
    private accountsById = new ObservableMap<number, AccountModel>();
    private loader: Loader;

    constructor(messageStore: IMessageStore, alertStore: AlertStore) {
        makeObservable(this);
        this.loader = new Loader(messageStore, alertStore);
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

    getAccount(id?: number) {
        return typeof id === 'number' ? this.accountsById.get(id) : undefined;
    }

    getCompany(id?: number) {
        return typeof id === 'number' ? this.companiesById.get(id) : undefined;
    }

    loadAccounts(): Promise<boolean> | undefined {
        if (!this.loading && this.accountsById.size === 0) {
            this.loading = true;
            return this.loader.load<IAccountsResponse>(loadingAccounts, {query,
                updater: ({accounts, companies}) => {
                    addToMap(this.companiesById, companies.map((company) => new CompanyModel(company)));
                    addToMap(this.accountsById, accounts.map((account) => new AccountModel(account, this.getCompany(account.companyId))));
                    for (const account of this.accounts) {
                        account.company?.accounts.push(account);
                    }
                },
                completer: () => this.loading = false,
            });
        }
    }

    saveCompanies(variables: IUpdateCompanies) {
        return this.loader.load<{companies: ICompany[]}>(savingCompanies, {query: updateCompaniesQuery, variables, updater: ({companies}) => {
            variables.delete.forEach(({id}) => this.companiesById.delete(id));
            companies.forEach((c) => {
                if (this.companiesById.has(c.id)) this.companiesById.get(c.id)?.update(c);
                else this.companiesById.set(c.id, new CompanyModel(c));
            });
        }});
    }
}
