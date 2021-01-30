import React from 'react';
import {shallow} from 'enzyme';
import TransactionModel from 'src/lib/model/TransactionModel';
import MessageStore from 'src/lib/store/MessageStore';
import PayeeStore from 'src/lib/store/PayeeStore';
import {PayeeModel} from 'src/lib/model/PayeeModel';
import Payee from './Payee';
import {newTx} from 'src/test/transactionFactory';
import CategoryStore from 'src/lib/store/CategoryStore';

const txData = newTx();

describe('Payee', () => {
    const messageStore = new MessageStore();
    const payeeStore = new PayeeStore(messageStore);
    const categoryStore = new CategoryStore(messageStore);

    beforeEach(() => {
        jest.spyOn(React, 'useContext').mockReturnValue({payeeStore});
    });
    it('returns null if no payee', () => {
        const tx = new TransactionModel(txData, categoryStore);

        const component = shallow(<Payee transaction={tx} />);

        expect(component).toBeEmptyRender();
    });
    it('shows payee', () => {
        const name = 'the payee';
        jest.spyOn(payeeStore, 'getPayee').mockReturnValue({name} as PayeeModel);
        const tx = new TransactionModel({...txData, payeeId: -1}, categoryStore);

        const component = shallow(<Payee transaction={tx} />);

        expect(component).toHaveClassName('payee');
        expect(component).toHaveClassName('chip');
        expect(component.find('i.material-icons.md-18')).toHaveText('person');
        expect(component.childAt(1)).toHaveText(name);
    });
});
