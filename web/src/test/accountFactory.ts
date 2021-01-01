import {AccountModel, IAccount, ICompany} from "src/lib/model/AccountModel";

let nextId = 0;

export function newCompany(overrides: Partial<ICompany> = {}): ICompany {
    return {
        id: `${++nextId}`,
        name: `Company ${nextId}`,
        version: 1,
        ...overrides,
    };
}

export function newAccount(overrides: Partial<IAccount> = {}, company?: ICompany): AccountModel {
    return new AccountModel({
        id: `${++nextId}`,
        name: `Account ${nextId}`,
        type: 'Bank',
        closed: false,
        version: 1,
        transactionCount: 555,
        balance: 1234.78,
        ...overrides,
    }, company);
}
