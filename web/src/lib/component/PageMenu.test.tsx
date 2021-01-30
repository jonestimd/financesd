import React from 'react';
import {shallow} from 'enzyme';
import PageMenu from './PageMenu';
import {MenuItem} from '@material-ui/core';
import {Link} from 'react-router-dom';

describe('PageMenu', () => {
    it('lists pages', () => {
        const component = shallow(<PageMenu />);

        const items = component.find(MenuItem);
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

        const items = component.find(MenuItem);
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
