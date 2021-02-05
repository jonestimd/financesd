import settingsStore from '../store/settingsStore';
import {CompanyModel} from './CompanyModel';
import {compareByName} from './entityUtils';

export interface IAccount {
    id: string;
    companyId?: number;
    type: string;
    name: string;
    description?: string;
    accountNo?: string;
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
    companyId?: number;
    type: string;
    name: string;
    description?: string;
    accountNo?: string;
    closed: boolean;
    version: number;
    transactionCount: number;
    balance: number;
    company?: CompanyModel;

    constructor(account: IAccount, company?: CompanyModel) {
        this.id = account.id;
        this.companyId = account.companyId;
        this.type = account.type;
        this.name = account.name;
        this.description = account.description;
        this.accountNo = account.accountNo;
        this.closed = account.closed;
        this.version = account.version;
        this.transactionCount = account.transactionCount;
        this.balance = account.balance;
        this.company = company;
    }

    get companyName(): string {
        return this.company ? this.company.name : '';
    }

    get displayName(): string {
        if (this.company) return `${this.company.name}: ${this.name}`;
        return this.name;
    }

    get hide() {
        return this.closed && settingsStore.hideClosedAccounts;
    }

    compareTo(that: AccountModel): number {
        const diff = compareByName(this.company, that.company);
        return diff === 0 ? compareByName(this, that) : diff;
    }
}
