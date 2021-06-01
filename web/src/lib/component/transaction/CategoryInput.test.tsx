import React from 'react';
import {shallow} from 'enzyme';
import CategoryInput, {IOption} from './CategoryInput';
import {newDetailModel} from 'test/detailFactory';
import Autocomplete, {AutocompleteRenderInputParams} from '@material-ui/lab/Autocomplete';
import {RootStore} from 'lib/store/RootStore';
import {newAccountModel} from 'test/accountFactory';
import {newCategoryModel} from 'test/categoryFactory';
import {Icon} from '@material-ui/core';

describe('CategoryInput', () => {
    const {accountStore, categoryStore} = new RootStore();
    const accounts = [newAccountModel(), newAccountModel()];
    accounts.forEach((a) => accountStore['accountsById'].set(a.id, a));
    const categories = [newCategoryModel(), newCategoryModel()];
    categories.forEach((c) => categoryStore['categoriesById'].set(c.id, c));
    const options = (categories as IOption[]).concat(accounts);

    beforeEach(() => {
        jest.spyOn(React, 'useContext').mockReturnValue({accountStore, categoryStore});
    });
    it('displays list of accounts and categories', () => {
        const detail = newDetailModel();

        const component = shallow(<CategoryInput detail={detail} />);

        expect(component.find(Autocomplete)).toHaveProp('options', options);
        expect(component.find(Autocomplete)).toHaveProp('value', null);
    });
    it('uses displayName for option labels', () => {
        const detail = newDetailModel();
        const component = shallow(<CategoryInput detail={detail} />);

        const getOptionLabel = component.find(Autocomplete).prop('getOptionLabel')!;

        expect(getOptionLabel(categories[0])).toEqual(categories[0].displayName);
        expect(getOptionLabel(accounts[0])).toEqual(accounts[0].displayName);
    });
    it('displays options containing text, ignoring case', () => {
        const detail = newDetailModel();
        const component = shallow(<CategoryInput detail={detail} />);

        const filterOptions = component.find(Autocomplete).prop('filterOptions')!;

        expect(filterOptions(options, {inputValue: 'CAT', getOptionLabel: () => ''})).toEqual(categories);
        expect(filterOptions(options, {inputValue: 'COUnt', getOptionLabel: () => ''})).toEqual(accounts);
    });
    it('displays selected category', () => {
        const detail = newDetailModel({transactionCategoryId: categories[1].id, categoryStore});

        const component = shallow(<CategoryInput detail={detail} />);

        expect(component.find(Autocomplete)).toHaveProp('value', categories[1]);
    });
    it('displays selected account', () => {
        const detail = newDetailModel({relatedDetail: {id: -2, transaction: {id: -1, accountId: accounts[1].id}}, accountStore});

        const component = shallow(<CategoryInput detail={detail} />);

        expect(component.find(Autocomplete)).toHaveProp('value', accounts[1]);
    });
    it('displays input with category icon', () => {
        const detail = newDetailModel();
        const component = shallow(<CategoryInput detail={detail} />);
        const renderInput = component.find(Autocomplete).prop('renderInput')!;

        const input = shallow(renderInput({} as AutocompleteRenderInputParams) as React.ReactElement);

        expect(input).toHaveProp('InputProps', expect.objectContaining({
            startAdornment: <Icon>category</Icon>,
        }));
    });
    it('sets detail.category', () => {
        const detail = newDetailModel();
        const component = shallow(<CategoryInput detail={detail} />);

        component.find(Autocomplete).simulate('change', {}, categories[0]);

        expect(detail.transactionCategoryId).toEqual(categories[0].id);
    });
});
