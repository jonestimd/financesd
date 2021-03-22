import {RootStore} from './RootStore';
import * as entityUtils from '../model/entityUtils';
import * as agent from '../agent';
import {AccountModel} from '../model/account/AccountModel';
import {newAccount, newAccountModel, newCompany, newCompanyModel} from 'src/test/accountFactory';
import {loadingAccounts, query, savingCompanies, updateCompaniesQuery} from './AccountStore';
import {CompanyModel} from '../model/account/CompanyModel';

describe('AccountStore', () => {
    const {accountStore, messageStore} = new RootStore();

    beforeEach(() => {
        accountStore['companiesById'].clear();
        accountStore['accountsById'].clear();
    });
    describe('get accounts', () => {
        it('sorts by company and name', () => {
            const account = newAccountModel();
            accountStore['accountsById'].set(account.id, account);
            jest.spyOn(entityUtils, 'sortValues');

            const accounts = accountStore.accounts;

            expect(accounts).toEqual([account]);
            expect(entityUtils.sortValues).toBeCalledWith(accountStore['accountsById'], AccountModel.compare);
        });
    });
    describe('get filteredAccounts', () => {
        it('returns accounts where hide is false', () => {
            const accounts = [newAccountModel(), newAccountModel()];
            accounts.forEach((account, index) => {
                jest.spyOn(account, 'hide', 'get').mockReturnValue(index === 0);
                accountStore['accountsById'].set(account.id, account);
            });

            expect(accountStore.filteredAccounts).toEqual(accounts.slice(1));
        });
    });
    describe('get accountsWithoutCompany', () => {
        it('returns accounts with no compnay where hide is false', () => {
            const accounts = [newAccountModel(), newAccountModel(), newAccountModel({}, newCompanyModel())];
            accounts.forEach((account, index) => {
                jest.spyOn(account, 'hide', 'get').mockReturnValue(index === 0);
                accountStore['accountsById'].set(account.id, account);
            });

            expect(accountStore.accountsWithoutCompany).toEqual(accounts.slice(1, 2));
        });
    });
    describe('get companies', () => {
        it('sorts by company name', () => {
            const company = newCompanyModel();
            accountStore['companiesById'].set(company.id, company);
            jest.spyOn(entityUtils, 'sortValuesByName');

            const companies = accountStore.companies;

            expect(companies).toEqual([company]);
            expect(entityUtils.sortValuesByName).toBeCalledWith(accountStore['companiesById']);
        });
    });
    describe('get filteredCompanies', () => {
        it('returns companies having unfiltered accounts', () => {
            const companies = [newCompanyModel(), newCompanyModel()];
            jest.spyOn(companies[0], 'filteredAccounts', 'get').mockReturnValue([newAccountModel()]);
            jest.spyOn(companies[1], 'filteredAccounts', 'get').mockReturnValue([]);
            companies.forEach((company) => accountStore['companiesById'].set(company.id, company));

            expect(accountStore.filteredCompanies).toEqual(companies.slice(0, 1));
        });
    });
    describe('getAccount', () => {
        it('returns account for ID', () => {
            const account = newAccountModel();
            accountStore['accountsById'].set(account.id, account);

            expect(accountStore.getAccount(account.id)).toBe(account);
            expect(accountStore.getAccount(parseInt(account.id))).toBe(account);
        });
        it('returns undefined for unknown ID', () => {
            expect(accountStore.getAccount('23')).toBeUndefined();
        });
    });
    describe('loadAccounts', () => {
        beforeEach(() => {
            accountStore['loading'] = false;
            jest.spyOn(messageStore, 'addProgressMessage');
            jest.spyOn(messageStore, 'removeProgressMessage');
        });
        it('loads accounts and companies if account map is empty', async () => {
            const company = newCompany();
            const account = newAccount();
            jest.spyOn(agent, 'graphql').mockResolvedValue({data: {companies: [company], accounts: [account]}});

            await accountStore.loadAccounts();

            expect(accountStore['loading']).toBe(false);
            expect(messageStore.addProgressMessage).toBeCalledWith(loadingAccounts);
            expect(messageStore.removeProgressMessage).toBeCalledWith(loadingAccounts);
            expect(agent.graphql).toBeCalledWith(query);
            expect(accountStore.companies).toEqual([new CompanyModel(company)]);
            expect(accountStore.accounts).toStrictEqual([new AccountModel(account)]);
        });
        it('does nothing is already loading', async () => {
            accountStore['loading'] = true;
            jest.spyOn(agent, 'graphql').mockRejectedValue(new Error());

            await accountStore.loadAccounts();

            expect(accountStore['loading']).toBe(true);
            expect(messageStore.addProgressMessage).not.toBeCalled();
            expect(messageStore.removeProgressMessage).not.toBeCalled();
            expect(agent.graphql).not.toBeCalled();
        });
        it('does nothing is already loaded', async () => {
            const account = newAccountModel();
            accountStore['accountsById'].set(account.id, account);
            jest.spyOn(agent, 'graphql').mockRejectedValue(new Error());

            await accountStore.loadAccounts();

            expect(accountStore['loading']).toBe(false);
            expect(messageStore.addProgressMessage).not.toBeCalled();
            expect(messageStore.removeProgressMessage).not.toBeCalled();
            expect(agent.graphql).not.toBeCalled();
        });
        it('logs error from graphql', async () => {
            const error = new Error('API error');
            jest.spyOn(agent, 'graphql').mockRejectedValue(error);
            jest.spyOn(console, 'error').mockImplementation(() => { });

            await accountStore.loadAccounts();

            expect(accountStore['loading']).toBe(false);
            expect(messageStore.addProgressMessage).toBeCalledWith(loadingAccounts);
            expect(messageStore.removeProgressMessage).toBeCalledWith(loadingAccounts);
            expect(console.error).toBeCalledWith('error gettting accounts', error);
        });
    });
    describe('saveCompanies', () => {
        const company1 = newCompanyModel();
        const company2 = newCompanyModel();

        beforeEach(() => {
            accountStore['companiesById'].set(company1.id, company1);
            accountStore['companiesById'].set(company2.id, company2);
            jest.spyOn(messageStore, 'addProgressMessage').mockReturnValue();
            jest.spyOn(messageStore, 'removeProgressMessage').mockReturnValue();
        });
        it('calls updateCompanies and adds new company', async () => {
            const changes = {add: ['new name'], update: [], delete: []};
            jest.spyOn(agent, 'graphql').mockResolvedValue({data: {companies: [{id: '-9', name: 'new name', version: 0}]}});

            expect(await accountStore.saveCompanies(changes)).toBe(true);

            expect(agent.graphql).toBeCalledWith(updateCompaniesQuery, changes);
            expect(messageStore.addProgressMessage).toBeCalledWith(savingCompanies);
            expect(messageStore.removeProgressMessage).toBeCalledWith(savingCompanies);
            expect(accountStore['companiesById'].get('-9')).toEqual(expect.objectContaining({name: 'new name', version: 0}));
        });
        it('calls updateCompanies and updates company', async () => {
            const changes = {add: [], update: [{id: company1.id, name: 'rename', version: company1.version}], delete: []};
            jest.spyOn(agent, 'graphql').mockResolvedValue({data: {companies: [
                {id: company1.id, name: 'rename', version: 2},
            ]}});

            expect(await accountStore.saveCompanies(changes)).toBe(true);

            expect(agent.graphql).toBeCalledWith(updateCompaniesQuery, changes);
            expect(messageStore.addProgressMessage).toBeCalledWith(savingCompanies);
            expect(messageStore.removeProgressMessage).toBeCalledWith(savingCompanies);
            expect(accountStore['companiesById'].get(company1.id)).toEqual(expect.objectContaining({name: 'rename', version: 2}));
        });
        it('calls updateCompanies and removes company', async () => {
            const changes = {add: [], update: [], delete: [parseInt(company2.id)]};
            jest.spyOn(agent, 'graphql').mockResolvedValue({data: {companies: []}});

            expect(await accountStore.saveCompanies(changes)).toBe(true);

            expect(agent.graphql).toBeCalledWith(updateCompaniesQuery, changes);
            expect(messageStore.addProgressMessage).toBeCalledWith(savingCompanies);
            expect(messageStore.removeProgressMessage).toBeCalledWith(savingCompanies);
            expect(accountStore['companiesById'].has(company2.id)).toBe(false);
        });
    });
});
