import React from 'react';
import {shallow} from 'enzyme';
import Category from './Category';
import MessageStore from 'src/lib/store/MessageStore';
import CategoryStore from 'src/lib/store/CategoryStore';
import {CategoryModel} from 'src/lib/model/CategoryModel';
import AccountStore from 'src/lib/store/AccountStore';
import {AccountModel} from 'src/lib/model/AccountModel';
import {newDetail} from 'src/test/detailFactory';

const detail = newDetail();

describe('Category', () => {
    const messageStore = new MessageStore();
    const accountStore = new AccountStore(messageStore);
    const categoryStore = new CategoryStore(messageStore);

    beforeEach(() => {
        jest.spyOn(React, 'useContext').mockReturnValue({accountStore, categoryStore});
    });
    it('returns null if no category', () => {
        const component = shallow(<Category detail={detail} />);

        expect(component).toBeEmptyRender();
    });
    it('shows category', () => {
        const transactionCategoryId = 123;
        const displayName = 'the category';
        jest.spyOn(categoryStore, 'getCategory').mockReturnValue({displayName} as CategoryModel);

        const component = shallow(<Category detail={{...detail, transactionCategoryId}} />);

        expect(component.find('.category')).toHaveText(displayName);
        expect(categoryStore.getCategory).toBeCalledWith(transactionCategoryId);
    });
    it('shows transfer account', () => {
        const name = 'the other account';
        const relatedDetail = {
            id: '2',
            amount: -detail.amount,
            transaction: {
                id: '20',
                accountId: 456,
            },
        };
        jest.spyOn(accountStore, 'getAccount').mockReturnValue({name} as AccountModel);

        const component = shallow(<Category detail={{...detail, relatedDetail}} />);

        expect(component.find('i.material-icons.md-18')).toHaveText('forward');
        expect(component.find('.transfer').childAt(1)).toHaveText(name);
        expect(accountStore.getAccount).toBeCalledWith(relatedDetail.transaction.accountId);
    });
});
