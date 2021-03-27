import CompanyRow from './CompanyRow';
import {newAccountModel} from 'src/test/accountFactory';
import {testAction} from 'src/test/mobxUtils';

describe('CompanyRow', () => {
    describe('get id', () => {
        it('returns data.id', () => {
            expect(new CompanyRow({id: -1}).id).toEqual('-1');
        });
    });
    describe('get version', () => {
        it('returns data version', () => {
            expect(new CompanyRow({id: -1, version: 1}).version).toEqual(1);
        });
        it('defaults to 0', () => {
            expect(new CompanyRow({id: -1}).version).toEqual(0);
        });
    });
    describe('get name', () => {
        it('returns data.name if no new name', () => {
            expect(new CompanyRow({id: -1, name: 'first name'}).name).toEqual('first name');
        });
        it('returns empty string if no data.name and no new name', () => {
            expect(new CompanyRow({id: -1}).name).toEqual('');
        });
        it('returns new name', () => {
            const model = new CompanyRow({id: -1});
            model.setName('new name');

            expect(model.name).toEqual('new name');
        });
    });
    describe('setName', () => {
        it('trims new name', () => {
            const model = new CompanyRow({id: -1});

            testAction(() => model.name, () => model.setName('  new name  '));

            expect(model.name).toEqual('new name');
            expect(model.isChanged).toBe(true);
        });
        it('undoes name change', () => {
            const model = new CompanyRow({id: -1, version: 1, name: 'old name'});
            model.setName('new name');

            testAction(() => model.name, () => model.setName('old name'));

            expect(model.name).toEqual('old name');
            expect(model.isChanged).toBe(false);
        });
    });
    describe('get accounts', () => {
        it('defaults to 0', () => {
            expect(new CompanyRow({id: -1}).accounts).toEqual(0);
        });
        it('returns accounts length', () => {
            expect(new CompanyRow({id: -1, accounts: [newAccountModel(), newAccountModel()]}).accounts).toEqual(2);
        });
    });
    describe('isChanged', () => {
        it('returns true if no version', () => {
            expect(new CompanyRow({id: -1}).isChanged).toBe(true);
        });
        it('returns true if new name', () => {
            const model = new CompanyRow({id: -1, version: 1});
            model.setName('new name');

            expect(model.isChanged).toBe(true);
        });
        it('returns false if version and no new name', () => {
            const model = new CompanyRow({id: -1, version: 1, name: 'old name'});

            expect(model.isChanged).toBe(false);
        });
    });
    describe('isValid', () => {
        it('returns true if name is not emtpy', () => {
            const model = new CompanyRow({id: -1, version: 1, name: 'old name'});
            model.setName('new name');

            expect(model.isValid).toBe(true);
        });
        it('returns false if name is emtpy', () => {
            const model = new CompanyRow({id: -1, version: 1, name: 'old name'});
            model.setName('');

            expect(model.isValid).toBe(false);
        });
    });
});
