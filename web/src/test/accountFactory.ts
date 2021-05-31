import {AccountModel, IAccount} from 'lib/model/account/AccountModel';
import {CompanyModel, ICompany} from 'lib/model/account/CompanyModel';

let nextId = 0;

export function newCompany(overrides: Partial<ICompany> = {}): ICompany {
    return {
        id: ++nextId,
        name: `Company ${nextId}`,
        version: 1,
        ...overrides,
    };
}

export function newCompanyModel(overrides: Partial<ICompany> = {}, ...accounts: AccountModel[]) {
    const company = new CompanyModel(newCompany(overrides), accounts);
    accounts.forEach((account) => account.company = company);
    return company;
}

export function newAccount(overrides: Partial<IAccount> = {}): IAccount {
    return {
        id: ++nextId,
        name: `Account ${nextId}`,
        type: 'Bank',
        closed: false,
        version: 1,
        transactionCount: 555,
        balance: 1234.78,
        ...overrides,
    };
}

export function newAccountModel(overrides: Partial<IAccount> = {}, company?: CompanyModel) {
    return new AccountModel(newAccount({...overrides, companyId: company && company.id}), company);
}
