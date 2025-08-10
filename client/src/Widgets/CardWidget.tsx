import React, {useEffect, useState} from 'react';
import {Button, Card, Tag} from 'antd';
import {DeleteTwoTone} from '@ant-design/icons';

import styled from 'styled-components';

const StyledCard = styled(Card)`
    position: relative;
    text-align: left;

    .ant-card-body {
        padding: 4px !important;
        margin-left: 0 !important;
    }

    .ant-card-head {
        margin-left: 0 !important;
        padding-top: 0px !important;
        padding-bottom: 0px !important;
    }

    .delete-button {
        position: absolute;
        top: 10px;
        right: 10px;
        background-color: transparent;
        border: none;
        cursor: pointer;
    }
`;

// Стилизованный компонент для уменьшенных тегов
const SmallTag = styled(Tag)`
    font-size: 10px;
    padding: 2px 2px;
    margin: 4px;
`;


const CardWidget = ({purchaseProp, onDelete}) => {
    const [purchase, setPurchase] = useState(purchaseProp)


    useEffect(() => {
        setPurchase(purchaseProp)
    }, [purchaseProp]);


    const handleDelete = (product_id) => {
        onDelete(product_id)
    };


    return (

        <StyledCard
            size="small"
            title={purchase.product.name}
            style={{
                boxShadow: '0 2px 8px rgba(0, 0, 0, 0.1)',
                padding: '0px',
                margin: '0 5px',
                position: 'relative',
                textAlign: 'left',
            }}
            extra={<Button type="secondary" shape="circle" style={{padding: '0px'}} icon={<DeleteTwoTone/>}
                           onClick={() => handleDelete(purchase.product.id)}/>}
        >

            <p style={{margin: '0px'}}>
                {purchase.price}₽ x {purchase.quantity} = {purchase.price * purchase.quantity} ₽
            </p>
            <span style={{marginBottom: '4px', paddingRight: '4px'}}>
                            <Tag color="green">{purchase.product.brand}</Tag>
                            <Tag color="blue">{purchase.product.volume}</Tag>
                            <Tag color="yellow">{purchase.product.category}</Tag>
                {purchase.tags && purchase.tags.map((tag, tagIndex) => (<SmallTag key={tagIndex} color="default">
                    {tag}
                </SmallTag>))}
                        </span>


        </StyledCard>);
};

export default CardWidget;
