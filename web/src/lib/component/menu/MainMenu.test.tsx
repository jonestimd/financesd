import React from 'react';
import {shallow, ShallowWrapper} from 'enzyme';
import MainMenu from './MainMenu';
import {IconButton, ListItemText, MenuItem} from '@material-ui/core';

const getMenuItems = (component: ShallowWrapper) => {
    return component.find(MenuItem).map((item) => {
        const listItemText = item.find(ListItemText);
        return [item.prop('to'), listItemText.length ? listItemText.text() : item.text()];
    });
};

describe('MainMenu', () => {
    const setDepth = jest.fn();

    describe('accounts page', () => {
        it('displays links to other pages', () => {
            const component = shallow(<MainMenu currentPage='menu.accounts' setDepth={setDepth} />);

            expect(getMenuItems(component)).toEqual([
                ['/finances/categories', 'Categories'],
                ['/finances/groups', 'Groups'],
                ['/finances/payees', 'Payees'],
                ['/finances/securities', 'Securities'],
            ]);
        });
    });
    const pages = [
        ['/finances/', 'Accounts'],
        ['/finances/categories', 'Categories'],
        ['/finances/groups', 'Groups'],
        ['/finances/payees', 'Payees'],
        ['/finances/securities', 'Securities'],
    ];
    [
        {name: 'categoriesPage', currentPage: 'menu.categories'},
        {name: 'groupsPage', currentPage: 'menu.groups'},
        {name: 'payeesPage', currentPage: 'menu.payees'},
        {name: 'securitiesPage', currentPage: 'menu.securities'},
        {name: 'transactionsPage', currentPage: 'transactions'},
    ].forEach(({name, currentPage}, index) => {
        describe(name, () => {
            it('displays links to other pages', () => {
                const component = shallow(<MainMenu currentPage={currentPage} setDepth={setDepth} />);

                const expectedPages = [...pages];
                expectedPages.splice(index + 1, 1);
                expect(getMenuItems(component)).toEqual(expectedPages);
            });
            it('displays button for companies menu', () => {
                const component = shallow(<MainMenu currentPage='menu.categories' setDepth={setDepth} />);

                component.find(IconButton).simulate('click');

                expect(setDepth).toBeCalledWith(1);
            });
        });
    });
});
