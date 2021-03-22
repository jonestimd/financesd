import {newAccount, newAccountModel, newCompanyModel} from 'src/test/accountFactory';
import settingsStore from '../../store/settingsStore';
import {AccountModel} from './AccountModel';

describe('AccountModel', () => {
    const account = newAccount();
    const company = newCompanyModel();

    describe('constructor', () => {
        it('populates account properties', () => {
            const model = new AccountModel(account);

            expect(model).toEqual(account);
        });
        it('sets company', () => {
            const model = new AccountModel(account, company);

            expect(model.company).toBe(company);
        });
    });
    describe('companyName', () => {
        it('returns empty string for no company', () => {
            const model = new AccountModel(account);

            expect(model.companyName).toEqual('');
        });
        it('returns company name', () => {
            const model = new AccountModel(account, company);

            expect(model.companyName).toEqual(company.name);
        });
    });
    describe('displayName', () => {
        it('returns name if no company', () => {
            const model = new AccountModel(account);

            expect(model.displayName).toEqual(account.name);
        });
        it('returns company name and name', () => {
            const model = new AccountModel(account, company);

            expect(model.displayName).toEqual(`${company.name}: ${account.name}`);
        });
    });
    describe('hide', () => {
        it('returns false if hide closed accounts is false', () => {
            jest.spyOn(settingsStore, 'hideClosedAccounts', 'get').mockReturnValue(false);

            expect(newAccountModel({closed: true}).hide).toEqual(false);
        });
        it('returns true if account is closed and hide closed accounts is true', () => {
            jest.spyOn(settingsStore, 'hideClosedAccounts', 'get').mockReturnValue(true);

            expect(newAccountModel().hide).toEqual(false);
            expect(newAccountModel({closed: true}).hide).toEqual(true);
        });
    });
    describe('compare', () => {
        it('sorts by company name then account name', () => {
            const account1 = newAccountModel({}, company);
            const account2 = newAccountModel({}, company);
            const company2 = newCompanyModel();
            const account3 = newAccountModel({name: account1.name}, company2);

            expect(AccountModel.compare(account1, account3)).toBeLessThan(0);
            expect(AccountModel.compare(account3, account1)).toBeGreaterThan(0);
            expect(AccountModel.compare(account1, account1)).toEqual(0);
            expect(AccountModel.compare(account1, account2)).toBeLessThan(0);
            expect(AccountModel.compare(account2, account1)).toBeGreaterThan(0);
        });
    });
});
