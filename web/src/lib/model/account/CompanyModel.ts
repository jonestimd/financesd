import {computed, makeObservable, observable} from 'mobx';
import {AccountModel} from './AccountModel';
import {compareByName} from '../entityUtils';

export interface ICompany {
    id: number;
    name: string;
    version: number;
}

export class CompanyModel implements ICompany {
    readonly id: number;
    @observable name: string;
    @observable version: number;
    @observable readonly accounts: AccountModel[];

    constructor(company: ICompany, accounts: AccountModel[] = []) {
        this.id = company.id;
        this.name = company.name;
        this.version = company.version;
        this.accounts = accounts;
        makeObservable(this);
    }

    @computed
    get filteredAccounts() {
        return this.accounts.filter((account) => !account.hide).sort(compareByName);
    }

    update(company: ICompany) {
        this.name = company.name;
        this.version = company.version;
    }
}
