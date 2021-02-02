import React from 'react';
import {shallow} from 'enzyme';
import PageMenu from './PageMenu';
import {MenuItem, MenuList} from '@material-ui/core';
import {Link} from 'react-router-dom';
import {RootStore} from '../store/RootStore';
import {newAccountModel, newCompanyModel} from 'src/test/accountFactory';

describe('PageMenu', () => {
    const rootStore = new RootStore();
    beforeEach(() => {
        jest.spyOn(React, 'useContext').mockReturnValue(rootStore);
    });
    describe('main menu', () => {
        it('lists pages', () => {
            const component = shallow(<PageMenu />);

            const items = component.find(MenuList).at(0).find(MenuItem);
            expect(items.map((item) => item.find(Link).prop('to'))).toEqual([
                '/finances/',
                '/finances/categories',
                '/finances/groups',
                '/finances/payees',
                '/finances/securities',
            ]);
            expect(items.map((item) => item.find(Link).text())).toEqual([
                'Accounts',
                'Categories',
                'Groups',
                'Payees',
                'Securities',
            ]);
        });
        it('excludes current page', () => {
            const component = shallow(<PageMenu currentPage='menu.categories' />);

            const items = component.find(MenuList).at(0).find(MenuItem);
            expect(items.map((item) => item.find(Link).prop('to'))).toEqual([
                '/finances/',
                '/finances/groups',
                '/finances/payees',
                '/finances/securities',
            ]);
            expect(items.map((item) => item.find(Link).text())).toEqual([
                'Accounts',
                'Groups',
                'Payees',
                'Securities',
            ]);
        });
    });
    describe('companies menu', () => {
        const {accountStore} = rootStore;
        it('lists accounts with no company', () => {
            const account = newAccountModel();
            jest.spyOn(accountStore, 'accounts', 'get').mockReturnValue([account]);

            const component = shallow(<PageMenu currentPage='menu.categories' />);

            const items = component.find(MenuList).at(1).find('MenuAccount').dive();
            expect(items.find(Link)).toHaveProps({to: `/finances/account/${account.id}`, children: account.name});
            expect(component.find(MenuList).at(1).find('MenuCompany')).toHaveLength(0);
        });
        it('lists accounts for selected company', () => {
            const company = newCompanyModel({}, newAccountModel(), newAccountModel());
            jest.spyOn(accountStore, 'companies', 'get').mockReturnValue([company]);
            jest.spyOn(accountStore, 'accounts', 'get').mockReturnValue(company.accounts);

            const component = shallow(<PageMenu currentPage='menu.categories' />);
            component.find(MenuList).at(1).find('MenuCompany').dive().find(MenuItem).simulate('click');

            const items = component.find(MenuList).at(2).find('MenuAccount');
            items.forEach((item, index) => {
                expect(item.dive().find(Link)).toHaveProps({
                    to: `/finances/account/${company.accounts[index].id}`,
                    children: company.accounts[index].name,
                });
            });
        });
        it('lists companies with a single account', () => {
            const company = newCompanyModel({}, newAccountModel());
            jest.spyOn(accountStore, 'companies', 'get').mockReturnValue([company]);

            const component = shallow(<PageMenu currentPage='menu.categories' />);

            expect(component.find(MenuList).at(1).find('MenuCompany').dive().dive().find(MenuItem))
                .toHaveText(`${company.name} \u25ba ${company.accounts[0].name}`);
        });
    });
});
