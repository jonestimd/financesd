import ListModel from '../ListModel';
import CompanyRow from './CompanyRow';
import {IOptions} from '../SelectionModel';
import {action} from 'mobx';
import AccountStore from '../../store/AccountStore';

export const nameKey = 'company.name';

export default class CompanyListModel extends ListModel<CompanyRow> {
    readonly accountStore: AccountStore;

    constructor(options: Omit<IOptions, 'rows'>, accountStore: AccountStore) {
        super(options, accountStore.companies.map((c) => new CompanyRow(c)));
        this.accountStore = accountStore;
    }

    @action
    addCompany() {
        this.add((id) => new CompanyRow({id: -id}));
    }

    validate = (row: CompanyRow, key: string) => {
        const errors: string[] = [];
        if (key === nameKey) {
            if (row.name.length === 0) errors.push('Name is required'); // TODO localization
            const name = row.name.toLowerCase();
            if (this.items.filter((c) => c.name.toLowerCase() === name).length > 1) errors.push('Name must be unique');
        }
        return errors;
    };

    async save() {
        const changes = {
            add: this.pendingAdds.map((c) => c.name),
            delete: Array.from(this.pendingDeletes).map((c) => c.id),
            update: this.changes.map(({id, name, version}) => ({id, name, version})),
        };
        const success = await this.accountStore.saveCompanies(changes);
        if (success) this.commit(this.accountStore.companies.map((c) => new CompanyRow(c)));
        return success;
    }
}
