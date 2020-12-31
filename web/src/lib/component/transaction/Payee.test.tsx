import React from 'react';
import Payee from './Payee';
import {shallow} from 'enzyme';
import TransactionModel, {ITransaction} from 'src/lib/model/TransactionModel';
import CategoryStore from 'src/lib/store/CategoryStore';
import MessageStore from 'src/lib/store/MessageStore';
import {CategoryModel} from 'src/lib/model/CategoryModel';
import PayeeStore from 'src/lib/store/PayeeStore';
import {PayeeModel} from 'src/lib/model/PayeeModel';

const txData: ITransaction = {
    id: '1',
    date: '2020-01-01',
    details: [],
    cleared: false,
};

describe('Payee', () => {
    const messageStore = new MessageStore();
    const categoryStore = new CategoryStore(messageStore);
    const payeeStore = new PayeeStore(messageStore);

    beforeEach(() => {
        jest.spyOn(React, 'useContext').mockReturnValue({payeeStore});
        jest.spyOn(categoryStore, 'getCategory').mockReturnValue({isAssetValue: false} as CategoryModel);
    });
    it('returns null if no payee', () => {
        const tx = new TransactionModel(txData, null);

        const component = shallow(<Payee transaction={tx} />);

        expect(component).toBeEmptyRender();
    });
    it('shows payee', () => {
        const name = 'the payee';
        jest.spyOn(payeeStore, 'getPayee').mockReturnValue({name} as PayeeModel);
        const tx = new TransactionModel({...txData, payeeId: -1}, null);

        const component = shallow(<Payee transaction={tx} />);

        expect(component).toHaveClassName('payee');
        expect(component).toHaveClassName('chip');
        expect(component.find('i.material-icons.md-18')).toHaveText('person');
        expect(component.childAt(1)).toHaveText(name);
    });
});
