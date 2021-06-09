import React from 'react';
import {Currency, Shares} from '../../formats';
import Category from './Category';
import Group from './Group';
import TransactionModel, {DetailField, securityDetailFields} from '../../model/TransactionModel';
import {observer} from 'mobx-react-lite';
import {Icon, IconButton, TextFieldProps} from '@material-ui/core';
import CategoryInput from './CategoryInput';
import classNames from 'classnames';
import GroupInput from './GroupInput';
import IconInput from '../IconInput';
import DetailModel from 'lib/model/DetailModel';
import NumberInput from '../NumberInput';

interface IProps {
    transaction: TransactionModel;
    detail: DetailModel;
    editField: DetailField | false;
    // TODO use to filter categories
    showSecurity?: boolean;
}

const fieldRenderers: Record<DetailField, (detail: DetailModel) => React.ReactNode> = {
    amount: ({amount, isValid}) => <span key='amount' className={classNames({error: !isValid})}><Currency>{amount}</Currency></span>,
    category: (detail) => <Category key='category' detail={detail} />,
    group: ({transactionGroupId}) => <Group key='group' id={transactionGroupId} />,
    shares: ({assetQuantity}) => assetQuantity ? <span key='shares' className='shares'><Icon>request_page</Icon><Shares>{assetQuantity}</Shares></span> : null,
    memo: ({memo}) => memo ? <span key='memo' className='memo'><Icon>notes</Icon>{memo}</span> : null,
};

const leading = ({transaction, detail, editField, deleted}: IProps & {deleted: boolean}) => {
    if (editField === securityDetailFields[0]) return null;
    const fields = editField ? securityDetailFields.slice(0, securityDetailFields.indexOf(editField)) : securityDetailFields;
    return (
        <div className={classNames('detail chip', {prefix: editField, deleted})} data-id={`${detail.id}`}>
            {!deleted && transaction.detailCount > 1
                ? <IconButton size='small' onClick={() => transaction.deleteDetail(detail)}><Icon>delete</Icon></IconButton>
                : null}
            {deleted
                ? <IconButton size='small' onClick={() => transaction.undeleteDetail(detail)}><Icon>undo</Icon></IconButton>
                : null}
            {fields.map((field: DetailField) => fieldRenderers[field](detail))}
        </div>
    );
};

const trailing = ({detail, editField, deleted}: Pick<IProps, 'detail' | 'editField'> & {deleted: boolean}) => {
    if (!editField || editField === securityDetailFields[securityDetailFields.length-1]) return null;
    const fields = securityDetailFields.slice(securityDetailFields.indexOf(editField)+1);
    return (
        <div className={classNames('detail chip suffix', {deleted})} data-id={`${detail.id}`}>
            {fields.map((field) => fieldRenderers[field](detail))}
            <span>&nbsp;</span>
        </div>
    );
};

// TODO validate number inputs: no alpha, precision, sign, shares required/not allowed
// TODO clean up amountText, assetQuantityText on blur

const TxDetail = observer<IProps & Partial<TextFieldProps>, HTMLDivElement>(({transaction, detail, editField, showSecurity, ...inputProps}) => {
    const {amountText, assetQuantityText, memo} = detail;
    const deleted = transaction.isDeleted(detail);
    return <>
        {leading({transaction, detail, editField, deleted})}
        {editField === 'amount' && <NumberInput {...inputProps} value={amountText} precision={2}
            onChange={(value) => detail.amountText = value} startAdornment={<span>$</span>} />}
        {editField === 'category' && <CategoryInput detail={detail} {...inputProps} />}
        {editField === 'group' && <GroupInput detail={detail} {...inputProps} />}
        {editField === 'shares' && <IconInput {...inputProps} type='number' value={assetQuantityText} icon='request_page'
            onChange={({currentTarget}) => detail.assetQuantityText = currentTarget.value} />}
        {editField === 'memo' && <IconInput {...inputProps} value={memo ?? ''} icon='notes'
            onChange={({currentTarget}) => detail.memo = currentTarget.value || undefined} />}
        {trailing({detail, editField, deleted})}
    </>;
}, {forwardRef: true});

export default TxDetail;
