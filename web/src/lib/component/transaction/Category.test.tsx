import React from 'react';
import {shallow} from 'enzyme';
import Category from './Category';
import {CategoryModel} from 'lib/model/CategoryModel';
import {AccountModel} from 'lib/model/account/AccountModel';
import {newDetailModel} from 'test/detailFactory';
import {RootStore} from 'lib/store/RootStore';
import {runInAction} from 'mobx';

const detail = newDetailModel();

describe('Category', () => {
    const rootStore = new RootStore();
    const {accountStore, categoryStore} = rootStore;

    beforeEach(() => {
        detail.reset();
        jest.spyOn(React, 'useContext').mockReturnValue(rootStore);
    });
    it('returns null if no category', () => {
        const component = shallow(<Category detail={detail} />);

        expect(component).toBeEmptyRender();
    });
    it('shows category', () => {
        const transactionCategoryId = 123;
        const displayName = 'the category';
        detail.transactionCategoryId = transactionCategoryId;
        jest.spyOn(categoryStore, 'getCategory').mockReturnValue({displayName} as CategoryModel);

        const component = shallow(<Category detail={detail} />);

        expect(component.find('.category')).toHaveText(displayName);
        expect(categoryStore.getCategory).toBeCalledWith(transactionCategoryId);
    });
    it('shows transfer account', () => {
        const name = 'the other account';
        const relatedDetail = {
            id: 2,
            amount: -detail.amount,
            transaction: {
                id: 20,
                accountId: 456,
            },
        };
        runInAction(() => detail.relatedDetail = relatedDetail);
        jest.spyOn(accountStore, 'getAccount').mockReturnValue({name} as AccountModel);

        const component = shallow(<Category detail={detail} />);

        expect(component.find('i.material-icons.md-18')).toHaveText('forward');
        expect(component.find('.transfer').childAt(1)).toHaveText(name);
        expect(accountStore.getAccount).toBeCalledWith(relatedDetail.transaction.accountId);
    });
});
