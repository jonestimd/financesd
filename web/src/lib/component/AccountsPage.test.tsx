import React from 'react';
import {shallow} from 'enzyme';
import {RootStore} from '../store/RootStore';
import {newAccountModel, newCompanyModel} from 'src/test/accountFactory';
import AccountsPage from './AccountsPage';
import Table from './table/Table';
import {IColumn} from './table/Row';
import TopAppBar from './TopAppBar';
import {AccountModel} from '../model/AccountModel';
import {Link} from 'react-router-dom';
import accountType from '../i18n/accountType';
import * as formats from '../formats';

describe('AccountsPage', () => {
    const rootStore = new RootStore();
    const account = newAccountModel({description: 'my special account', accountNo: '123456-789'}, newCompanyModel());

    beforeEach(() => {
        jest.spyOn(React, 'useContext').mockReturnValue(rootStore);
        jest.spyOn(rootStore.accountStore, 'loadAccounts').mockResolvedValue();
    });
    it('displays app bar and table of accounts', () => {
        const accounts = [account];
        jest.spyOn(rootStore.accountStore, 'accounts', 'get').mockReturnValue(accounts);

        const component = shallow(<AccountsPage />);

        expect(component.find(TopAppBar)).toHaveProp('title', 'Accounts');
        expect(component.find(Table)).toHaveProp('data', accounts);
    });
    describe('accounts table', () => {
        const columnTests = [
            {name: 'company', value: account.companyName},
            {name: 'name', value: <Link to={`account/${account.id}`}>{account.name}</Link>},
            {name: 'type', value: accountType(account.type)},
            {name: 'description', value: account.description},
            {name: 'number', value: account.accountNo},
            {name: 'closed', value: null},
            {name: 'transactions', value: account.transactionCount},
            {name: 'balance', value: formats.currency.format(account.balance)},
        ];

        columnTests.forEach(({name, value}) => {
            it(`displays ${name}`, () => {
                const component = shallow(<AccountsPage />);
                const columns = component.find(Table).prop<IColumn<AccountModel>[]>('columns');

                const column = columns.find(({key}) => key === `account.${name}`);

                expect(column?.render(account)).toEqual(value);
            });
        });
    });
});
