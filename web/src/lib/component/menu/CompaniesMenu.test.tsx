import React from 'react';
import {shallow} from 'enzyme';
import CompaniesMenu from './CompaniesMenu';
import {RootStore} from 'src/lib/store/RootStore';
import {newAccountModel, newCompanyModel} from 'src/test/accountFactory';
import {MenuAccount} from './AccountsMenu';
import settingsStore from 'src/lib/store/settingsStore';
import {Checkbox, ListItem, ListItemText, MenuItem} from '@material-ui/core';
import {runInAction} from 'mobx';
import ChildMenu from './ChildMenu';

describe('CompanyMenu', () => {
    const props = {
        onBack: jest.fn(),
        selectCompany: jest.fn(),
    };
    const rootStore = new RootStore();
    beforeEach(() => {
        runInAction(() => {
            rootStore.accountStore['accountsById'].clear();
            rootStore.accountStore['companiesById'].clear();
        });
        jest.spyOn(React, 'useContext').mockReturnValue(rootStore);
        jest.spyOn(settingsStore, 'hideClosedAccounts', 'get').mockReturnValue(false);
    });
    it('displays checkbox to hide closed accounts', () => {
        const setter = jest.spyOn(settingsStore, 'hideClosedAccounts', 'set');

        const component = shallow(<CompaniesMenu {...props} />);

        const listItem = component.find(ListItem);
        expect(listItem).toHaveText('Hide closed accounts');
        const checkbox = listItem.find(Checkbox);
        expect(checkbox).toHaveProp('checked', false);
        checkbox.simulate('change');
        expect(setter).toBeCalledWith(true);
    });
    it('displays accounts with no company', () => {
        const account = newAccountModel();
        runInAction(() => rootStore.accountStore['accountsById'].set(account.id, account));

        const component = shallow(<CompaniesMenu {...props} />);

        expect(component.find(MenuAccount)).toHaveProps({...account});
    });
    it('hides closed accounts with no company', () => {
        jest.spyOn(settingsStore, 'hideClosedAccounts', 'get').mockReturnValue(true);
        const account = newAccountModel({closed: true});
        runInAction(() => rootStore.accountStore['accountsById'].set(account.id, account));

        const component = shallow(<CompaniesMenu {...props} />);

        expect(component.find(MenuAccount)).not.toExist();
    });
    it('hides companies with no accounts', () => {
        const company = newCompanyModel();
        runInAction(() => rootStore.accountStore['companiesById'].set(company.id, company));

        const component = shallow(<CompaniesMenu {...props} />);

        expect(component.find(ChildMenu).prop<unknown[]>('children')[2]).toHaveLength(0);
    });
    it('displays account link for company with single account', () => {
        const company = newCompanyModel({}, newAccountModel());
        runInAction(() => rootStore.accountStore['companiesById'].set(company.id, company));

        const component = shallow(<CompaniesMenu {...props} />);

        expect(component.find('*').last().dive().find(MenuAccount)).toHaveProps({...company.accounts[0]});
    });
    it('displays menu item for company with multiple accounts', () => {
        const company = newCompanyModel({}, newAccountModel(), newAccountModel());
        runInAction(() => rootStore.accountStore['companiesById'].set(company.id, company));

        const component = shallow(<CompaniesMenu {...props} />);

        const menuItem = component.find('*').last().dive().find(MenuItem);
        expect(menuItem.find(ListItemText)).toHaveText(company.name);
        menuItem.simulate('click');
        expect(props.selectCompany).toBeCalledWith(company);
    });
});
