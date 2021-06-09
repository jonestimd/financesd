import {Typography, Checkbox, TextFieldProps} from '@material-ui/core';
import classNames from 'classnames';
import {observer} from 'mobx-react-lite';
import React from 'react';
import TransactionModel from '../../model/TransactionModel';
import Memo from './Memo';
import Payee from './Payee';
import Security from './Security';
import TxDetail from './TxDetail';
import * as formats from '../../formats';
import PayeeInput from './PayeeInput';
import SecurityInput from './SecurityInput';
import DateInput from '../DateInput';
import IconInput from '../IconInput';

interface IProps {
    tx: TransactionModel;
    selected: boolean;
    showSecurity?: boolean;
    fieldIndex: number;
    setField: (field: number) => void;
}

const unselectedField = {transactionField: false as const, detailField: false as const, itemIndex: -1};

const Transaction: React.FC<IProps> = observer(({tx, selected, showSecurity, fieldIndex, setField}) => {
    const {transactionField, detailField, itemIndex} = selected ? tx.getField(fieldIndex, showSecurity) : unselectedField;
    const inputProps: Pick<TextFieldProps, 'variant' | 'color' | 'size' | 'onKeyDown'> = {
        variant: 'outlined',
        color: 'primary',
        size: 'small',
        onKeyDown: (event: React.KeyboardEvent<HTMLInputElement>) => {
            const {key, ctrlKey, shiftKey, altKey} = event;
            if (key === 'Tab' && !ctrlKey && !altKey) {
                event.stopPropagation();
                event.preventDefault();
                setField(tx.clampField(fieldIndex + (shiftKey ? -1 : 1), showSecurity));
            }
        },
    };
    return (
        <Typography className={classNames('transaction', {selected, changed: tx.isChanged})}>
            <div className='leading'>
                {transactionField === 'date'
                    ? <DateInput initialValue={tx.date ?? ''} {...inputProps} onDateChange={(_date, value) => tx.date = value} />
                    : <span className='date'>{tx.date}</span>}
                {transactionField === 'ref'
                    ? <IconInput value={tx.referenceNumber ?? ''} {...inputProps} icon='tag' onChange={(event) => tx.referenceNumber = event.currentTarget.value} />
                    : tx.referenceNumber ? <span data-type='ref'>{tx.referenceNumber}</span> : null}
            </div>
            <div className='details'>
                {transactionField === 'payee' ? <PayeeInput transaction={tx} {...inputProps} /> : <Payee transaction={tx} />}
                {transactionField === 'security' ? <SecurityInput transaction={tx} {...inputProps} /> : <Security transaction={tx} />}
                {transactionField === 'description'
                    ? <IconInput value={tx.memo ?? ''} {...inputProps} icon='notes' onChange={(event) => tx.memo = event.currentTarget.value} />
                    : <Memo text={tx.memo} />}
                {tx.details.map((detail, index) =>
                    <TxDetail key={detail.id ?? -index} detail={detail} transaction={tx}
                        editField={index === itemIndex && detailField} showSecurity={showSecurity} {...inputProps} />)}
            </div>
            <div className='trailing'>
                <Checkbox checked={tx.cleared} onChange={(_event, checked) => tx.cleared = checked} />
                <span className={formats.numberClass(tx.subtotal)}>{formats.currency.format(tx.subtotal)}</span>
                <span className={formats.numberClass(tx.balance)}>{formats.currency.format(tx.balance)}</span>
            </div>
        </Typography>
    );
});

export default Transaction;
