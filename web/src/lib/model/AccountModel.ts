import {compareByName} from './entityUtils';

export interface ICompany {
    id: string;
    name: string;
    version: number;
}

export interface IAccount {
    id: string;
    companyId: number;
    type: string;
    name: string;
    description: string;
    accountNo: string;
    closed: boolean;
    version: number;
    transactionCount: number;
    balance: number;
}

export class AccountModel implements IAccount {
    static compare(a1: AccountModel, a2: AccountModel): number {
        return a1.compareTo(a2);
    }

    id: string;
    companyId: number;
    type: string;
    name: string;
    description: string;
    accountNo: string;
    closed: boolean;
    version: number;
    transactionCount: number;
    balance: number;
    company?: ICompany;

    constructor(account: IAccount, company?: ICompany) {
        Object.assign(this, account);
        this.company = company;
    }

    get companyName(): string {
        return this.company ? this.company.name : '';
    }

    get displayName(): string {
        if (this.company) return `${this.company.name}: ${this.name}`;
        return this.name;
    }

    compareTo(that: AccountModel): number {
        const diff = compareByName(this.company, that.company);
        return diff === 0 ? compareByName(this, that) : diff;
    }
}
