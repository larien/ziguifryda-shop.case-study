create database TrabSO;

use TrabSO;
----------------------------------------Tabelas----------------------------------------------
CREATE TABLE users(
	id int PRIMARY KEY AUTO_INCREMENT,
	login varchar(8) not null,
	password varchar(8) not null,
	name1 varchar(20) not null,
    name2 varchar(20),
	surname varchar (20),
	CPF varchar(12) not null,
	Email varchar(50) not null,
	Telefone integer not null,
    dataNasc date not null,
	FKEndereco int not null
);


CREATE TABLE Enderecos(
	id int PRIMARY KEY AUTO_INCREMENT,
	street varchar(200) not null,
    neighborhood varchar(200) not null,
	city varchar(200) not null,
	number int not null,
    complement varchar(200),
	state varchar(2) not null
);

Create TABLE orders(
IDPedido int  PRIMARY KEY AUTO_INCREMENT,
user_id int not null,
product_id integer not null,
IDFormaPagamento int not null,
IDEnderecoEntrega int not null,
price float not null,
sell_date date not null,
IdStatus int  not null

);

CREATE TABLE FormaPagamento(
	IDFormaPagamento int  PRIMARY KEY AUTO_INCREMENT,
	Forma_de_Pagamento varchar(200)
);
CREATE TABLE Status (
IdStatus int  PRIMARY KEY AUTO_INCREMENT,
situacao varchar(30)
);

CREATE TABLE products(
id int PRIMARY KEY AUTO_INCREMENT,
title varchar(30),
description varchar(50),
price float not null,
quantity integer not null default 0,
filename varchar(30)

);

CREATE TABLE ItensPedido(
IdItens int PRIMARY KEY Auto_INCREMENT,
IdProduto int not null,
IDPedido int not null,
Quantidade int not null
);

-----------------------------------------Chaves Estrangeiras------------------------------------------------------

ALTER TABLE users ADD CONSTRAINT FKClienteEndereco FOREIGN KEY (FKEndereco) REFERENCES Enderecos(id);
ALTER TABLE orders ADD CONSTRAINT FKPedidoCliente FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE orders ADD CONSTRAINT FKPedidoPagamento FOREIGN KEY (IDFormaPagamento) REFERENCES FormaPagamento(IDFormaPagamento);
ALTER TABLE orders ADD CONSTRAINT FKEnderecoEntrega FOREIGN KEY (IDEnderecoEntrega) REFERENCES Enderecos(id);
ALTER TABLE orders ADD CONSTRAINT FKStatus FOREIGN KEY (IdStatus) REFERENCES  Status(IdStatus);
ALTER TABLE ItensPedido ADD CONSTRAINT FKProduto FOREIGN KEY (IdProduto) REFERENCES products(id);
ALTER TABLE ItensPedido ADD CONSTRAINT FKPedido FOREIGN KEY (IDPedido) REFERENCES Pedido(IDPedido);
 
