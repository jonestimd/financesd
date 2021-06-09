import React from 'react';
import {shallow} from 'enzyme';
import TxDetail from './TxDetail';
import Category from './Category';
import {Currency, Shares} from 'lib/formats';
import Group from './Group';
import {newDetail, newDetailModel} from 'test/detailFactory';
import {Icon, IconButton} from '@material-ui/core';
import CategoryInput from './CategoryInput';
import GroupInput from './GroupInput';
import IconInput from '../IconInput';
import NumberInput from '../NumberInput';
import {newTxModel} from '../../../test/transactionFactory';

const tx1 = newTxModel({details: [newDetail({transactionGroupId: 456})]});
const detail = tx1.details[0];
const props = {detail, transaction: tx1};

describe('TxDetail', () => {
    const tx2 = newTxModel({details: [newDetail(), newDetail()]});

    beforeEach(() => {
        tx1.reset();
        tx2.reset();
    });
    it('shows amount, category and group', () => {
        const component = shallow(<TxDetail {...props} editField={false} />);

        expect(component.find(Currency)).toHaveProp('children', detail.amount);
        expect(component.find(Category)).toHaveProp('detail', detail);
        expect(component.find(Group)).toHaveProp('id', detail.transactionGroupId);
        expect(component.find('.shares')).not.toExist();
        expect(component.find('.memo')).not.toExist();
    });
    it('highlights invalid amount', () => {
        detail.amountText = '';

        const component = shallow(<TxDetail {...props} editField={false} />);

        expect(component.find('span').at(0)).toHaveClassName('error');
    });
    it('shows asset quantity', () => {
        const assetQuantity = 100.789;
        detail.assetQuantity = assetQuantity;

        const component = shallow(<TxDetail {...props} editField={false} />);

        expect(component.find('.shares').find(Shares)).toHaveProp('children', assetQuantity);
    });
    it('shows memo', () => {
        const memo = 'note to self';
        detail.memo = memo;

        const component = shallow(<TxDetail {...props} editField={false} />);

        expect(component.find('.memo')).toHaveProp('children', [<Icon>notes</Icon>, memo]);
    });
    it('edits amount', () => {
        const component = shallow(<TxDetail {...props} editField='amount' />);
        const value = detail.amountText;

        const input = component.find(NumberInput);
        input.simulate('change', '234.56');

        expect(input).toHaveProp({value, precision: 2, startAdornment: <span>$</span>});
        expect(detail.amount).toEqual(234.56);
        expect(detail.amountText).toEqual('234.56');
    });
    it('edits category', () => {
        const component = shallow(<TxDetail {...props} editField='category' />);

        const input = component.find(CategoryInput);

        expect(input).toHaveProp('detail', detail);
    });
    it('edits group', () => {
        const component = shallow(<TxDetail {...props} editField='group' />);

        const input = component.find(GroupInput);

        expect(input).toHaveProp('detail', detail);
    });
    it('edits shares', () => {
        const detail = newDetailModel({assetQuantity: 123});
        const component = shallow(<TxDetail {...props} detail={detail} editField='shares' />);
        const value = detail.assetQuantityText;

        const input = component.find(IconInput);
        input.simulate('change', {currentTarget: {value: '234.567'}});

        expect(input).toHaveProp({type: 'number', value, icon: 'request_page'});
        expect(detail.assetQuantity).toEqual(234.567);
        expect(detail.assetQuantityText).toEqual('234.567');
    });
    it('highlights deleted detail', () => {
        tx2.deleteDetail(tx2.details[0]);

        const component = shallow(<TxDetail transaction={tx2} detail={tx2.details[0]} editField='shares' />);

        expect(component.find('div.chip')).toHaveClassName('deleted');
    });
    describe('delete button', () => {
        it('adds detail to pending deletes', () => {
            const component = shallow(<TxDetail transaction={tx2} detail={tx2.details[0]} editField='shares' />);

            component.find(IconButton).simulate('click');

            expect(tx2.deletedDetails).toContain(tx2.details[0].id);
            expect(component.find(IconButton).find(Icon)).toHaveText('delete');
        });
        it('is not displayed for single detail', () => {
            const component = shallow(<TxDetail {...props} editField='shares' />);

            expect(component.find(IconButton)).not.toExist();
        });
        it('is not displayed for single undeleted detail', () => {
            tx2.deleteDetail(tx2.details[0]);

            const component = shallow(<TxDetail transaction={tx2} detail={tx2.details[1]} editField='shares' />);

            expect(component.find(IconButton)).not.toExist();
        });
    });
    describe('undo button', () => {
        it('removes detail from pending deletes', () => {
            tx2.deleteDetail(tx2.details[0]);
            const component = shallow(<TxDetail transaction={tx2} detail={tx2.details[0]} editField='shares' />);

            component.find(IconButton).simulate('click');

            expect(tx2.deletedDetails).toHaveLength(0);
            expect(component.find(IconButton).find(Icon)).toHaveText('undo');
        });
    });
    describe('memo input', () => {
        it('displays existing memo', () => {
            const detail = newDetailModel({memo: 'some notes'});
            const component = shallow(<TxDetail {...props} detail={detail} editField='memo' />);

            const input = component.find(IconInput);

            expect(input).toHaveProp({value: detail.memo, icon: 'notes'});
        });
        it('displays blank for no memo', () => {
            const detail = newDetailModel();
            const component = shallow(<TxDetail {...props} detail={detail} editField='memo' />);

            const input = component.find(IconInput);

            expect(input).toHaveProp({value: '', icon: 'notes'});
        });
        it('sets memo on detail', () => {
            const detail = newDetailModel();
            const component = shallow(<TxDetail {...props} detail={detail} editField='memo' />);

            component.find(IconInput).simulate('change', {currentTarget: {value: 'some notes'}});

            expect(detail.memo).toEqual('some notes');
        });
        it('clears memo on detail', () => {
            const detail = newDetailModel({memo: 'some notes'});
            const component = shallow(<TxDetail {...props} detail={detail} editField='memo' />);

            component.find(IconInput).simulate('change', {currentTarget: {value: ''}});

            expect(detail.memo).toBeUndefined();
        });
    });
});
