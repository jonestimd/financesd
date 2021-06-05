import React from 'react';
import {shallow} from 'enzyme';
import TxDetail from './TxDetail';
import Category from './Category';
import {Currency, Shares} from 'lib/formats';
import Group from './Group';
import {newDetailModel} from 'test/detailFactory';
import {Icon} from '@material-ui/core';
import CategoryInput from './CategoryInput';
import GroupInput from './GroupInput';
import IconInput from '../IconInput';
import NumberInput from '../NumberInput';

const detail = newDetailModel({transactionGroupId: 456});

describe('TxDetail', () => {
    beforeEach(() => detail.reset());
    it('shows amount, category and group', () => {
        const component = shallow(<TxDetail detail={detail} editField={false} />);

        expect(component.find(Currency)).toHaveProp('children', detail.amount);
        expect(component.find(Category)).toHaveProp('detail', detail);
        expect(component.find(Group)).toHaveProp('id', detail.transactionGroupId);
        expect(component.find('.shares')).not.toExist();
        expect(component.find('.memo')).not.toExist();
    });
    it('shows asset quantity', () => {
        const assetQuantity = 100.789;
        detail.assetQuantity = assetQuantity;

        const component = shallow(<TxDetail detail={detail} editField={false} />);

        expect(component.find('.shares').find(Shares)).toHaveProp('children', assetQuantity);
    });
    it('shows memo', () => {
        const memo = 'note to self';
        detail.memo = memo;

        const component = shallow(<TxDetail detail={detail} editField={false} />);

        expect(component.find('.memo')).toHaveProp('children', [<Icon>notes</Icon>, memo]);
    });
    it('edits amount', () => {
        const component = shallow(<TxDetail detail={detail} editField='amount' />);
        const value = detail.amountText;

        const input = component.find(NumberInput);
        input.simulate('change', '234.56');

        expect(input).toHaveProp({value, precision: 2, startAdornment: <span>$</span>});
        expect(detail.amount).toEqual(234.56);
        expect(detail.amountText).toEqual('234.56');
    });
    it('edits category', () => {
        const component = shallow(<TxDetail detail={detail} editField='category' />);

        const input = component.find(CategoryInput);

        expect(input).toHaveProp('detail', detail);
    });
    it('edits group', () => {
        const component = shallow(<TxDetail detail={detail} editField='group' />);

        const input = component.find(GroupInput);

        expect(input).toHaveProp('detail', detail);
    });
    it('edits shares', () => {
        const detail = newDetailModel({assetQuantity: 123});
        const component = shallow(<TxDetail detail={detail} editField='shares' />);
        const value = detail.assetQuantityText;

        const input = component.find(IconInput);
        input.simulate('change', {currentTarget: {value: '234.567'}});

        expect(input).toHaveProp({type: 'number', value, icon: 'request_page'});
        expect(detail.assetQuantity).toEqual(234.567);
        expect(detail.assetQuantityText).toEqual('234.567');
    });
    describe('memo input', () => {
        it('displays existing memo', () => {
            const detail = newDetailModel({memo: 'some notes'});
            const component = shallow(<TxDetail detail={detail} editField='memo' />);

            const input = component.find(IconInput);

            expect(input).toHaveProp({value: detail.memo, icon: 'notes'});
        });
        it('displays blank for no memo', () => {
            const detail = newDetailModel();
            const component = shallow(<TxDetail detail={detail} editField='memo' />);

            const input = component.find(IconInput);

            expect(input).toHaveProp({value: '', icon: 'notes'});
        });
        it('sets memo on detail', () => {
            const detail = newDetailModel();
            const component = shallow(<TxDetail detail={detail} editField='memo' />);

            component.find(IconInput).simulate('change', {currentTarget: {value: 'some notes'}});

            expect(detail.memo).toEqual('some notes');
        });
        it('clears memo on detail', () => {
            const detail = newDetailModel({memo: 'some notes'});
            const component = shallow(<TxDetail detail={detail} editField='memo' />);

            component.find(IconInput).simulate('change', {currentTarget: {value: ''}});

            expect(detail.memo).toBeUndefined();
        });
    });
});
