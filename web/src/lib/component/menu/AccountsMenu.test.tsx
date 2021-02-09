import React from 'react';
import {shallow} from 'enzyme';
import {newAccountModel, newCompanyModel} from 'src/test/accountFactory';
import AccountsMenu, {MenuAccount} from './AccountsMenu';
import {ListItemText, MenuItem} from '@material-ui/core';
import {Link} from 'react-router-dom';
import ChildMenu from './ChildMenu';

describe('AccountsMenu', () => {
    const onBack = jest.fn();
    it('displays back button', () => {
        const company = newCompanyModel({}, newAccountModel(), newAccountModel());

        const component = shallow(<AccountsMenu onBack={onBack} company={company} />);

        component.find(ChildMenu).prop('onBack')();
        expect(onBack).toBeCalled();
    });
    it('displays company name and accounts', () => {
        const company = newCompanyModel({}, newAccountModel(), newAccountModel());

        const component = shallow(<AccountsMenu onBack={onBack} company={company} />);

        expect(component.find(ChildMenu)).toHaveProp('name', company.name);
        const items = component.find(MenuAccount);
        items.forEach((item, index) => {
            const account = company.accounts[index];
            const menuItem = item.dive();
            expect(menuItem.find(MenuItem)).toHaveProps({component: Link, to: `/finances/account/${account.id}`});
            expect(menuItem.find(ListItemText)).toHaveProp('secondary', null);
            expect(menuItem.find(ListItemText)).toHaveText(account.name);
        });
    });
    it('hides closed accounts', () => {
        const company = newCompanyModel({}, newAccountModel(), newAccountModel(), newAccountModel());
        jest.spyOn(company.accounts[0], 'hide', 'get').mockReturnValue(true);

        const component = shallow(<AccountsMenu onBack={onBack} company={company} />);

        expect(component.find(ChildMenu)).toHaveProp('name', company.name);
        const items = component.find(MenuAccount);
        items.forEach((item, index) => {
            const account = company.accounts[index + 1];
            const menuItem = item.dive();
            expect(menuItem.find(MenuItem)).toHaveProps({component: Link, to: `/finances/account/${account.id}`});
            expect(menuItem.find(ListItemText)).toHaveProp('secondary', null);
            expect(menuItem.find(ListItemText)).toHaveText(account.name);
        });
    });
    describe('MenuAccount', () => {
        it('displays account with no company', () => {
            const account = newAccountModel();

            const component = shallow(<MenuAccount {...account} />);

            expect(component.find(MenuItem)).toHaveProps({component: Link, to: `/finances/account/${account.id}`});
            expect(component.find(ListItemText)).toHaveProp('secondary', null);
            expect(component.find(ListItemText)).toHaveText(account.name);
        });
        it('displays company name for company with single account', () => {
            const account = newAccountModel();
            const company = newCompanyModel({}, account);

            const component = shallow(<MenuAccount {...account} />);

            expect(component.find(MenuItem)).toHaveProps({component: Link, to: `/finances/account/${account.id}`});
            expect(component.find(ListItemText)).toHaveProp('secondary', account.name);
            expect(component.find(ListItemText)).toHaveText(company.name);
        });
    });
});
