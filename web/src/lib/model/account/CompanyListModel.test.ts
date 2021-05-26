import {RootStore} from 'src/lib/store/RootStore';
import {newCompanyModel} from 'src/test/accountFactory';
import {testAction} from 'src/test/mobxUtils';
import CompanyListModel, {nameKey} from './CompanyListModel';
import CompanyRow from './CompanyRow';

describe('CompanyListModel', () => {
    const {accountStore} = new RootStore();
    const companies = [newCompanyModel(), newCompanyModel()];

    beforeEach(() => {
        jest.spyOn(accountStore, 'companies', 'get').mockReturnValue(companies);
    });
    describe('constructor', () => {
        it('sets items', () => {
            const model = new CompanyListModel({columns: 2}, accountStore);

            expect(model.items).toEqual(companies.map((c) => new CompanyRow(c)));
        });
    });
    describe('addCompany', () => {
        it('adds a new company', () => {
            const model = new CompanyListModel({columns: 2}, accountStore);

            testAction(() => model.pendingAdds, () => model.addCompany());

            expect(model.pendingAdds).toHaveLength(1);
            expect(model.pendingAdds[0].id).toBeDefined();
        });
    });
    describe('validate', () => {
        it('returns empty array for unknown key', () => {
            const model = new CompanyListModel({columns: 2}, accountStore);

            expect(model.validate(new CompanyRow({id: -1}), '??')).toHaveLength(0);
        });
        it('returns empty array for valid name', () => {
            const model = new CompanyListModel({columns: 2}, accountStore);

            expect(model.validate(model.items[0], nameKey)).toHaveLength(0);
        });
        it('returns error for empty name', () => {
            const model = new CompanyListModel({columns: 2}, accountStore);

            expect(model.validate(new CompanyRow({id: -1}), nameKey)).toContain('Name is required');
        });
        it('returns error for duplicate name', () => {
            const model = new CompanyListModel({columns: 2}, accountStore);
            model.addCompany();
            const item = model.items[2];
            item.setName(companies[0].name);

            expect(model.validate(item, nameKey)).toContain('Name must be unique');
        });
    });
    describe('save', () => {
        const updatedCompanies = [newCompanyModel(), newCompanyModel()];
        const prepareModel = () => {
            const model = new CompanyListModel({columns: 2}, accountStore);
            model.delete();
            model.items[1].setName('new name');
            model.addCompany();
            model.pendingAdds[0].setName('new item');
            return model;
        };

        beforeEach(() => {
            jest.spyOn(accountStore, 'companies', 'get')
                .mockReturnValueOnce(companies)
                .mockReturnValueOnce(updatedCompanies);
        });
        it('commits on success', async () => {
            jest.spyOn(accountStore, 'saveCompanies').mockResolvedValue(true);
            const model = prepareModel();
            const add = model.pendingAdds.map((c) => c.name);
            const del = model.pendingDeletes.map(({id, version}) => ({id, version}));
            const update = model.changes.map(({id, name, version}) => ({id, name, version}));

            expect(await model.save()).toBe(true);

            expect(accountStore.saveCompanies).toBeCalledWith({add, update, delete: del});
            expect(model.isChanged).toBe(false);
            expect(model.items.length).toEqual(2);
            expect(model.items).toEqual(updatedCompanies.map((c) => new CompanyRow(c)));
        });
        it('does not commit on error', async () => {
            jest.spyOn(accountStore, 'saveCompanies').mockResolvedValue(false);
            const model = prepareModel();

            expect(await model.save()).toBe(false);

            expect(model.isChanged).toBe(true);
        });
    });
});
