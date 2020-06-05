# Sena Challenge

## Descrição do desafio

Receber uma sequencia de 6 numeros aleatórios entre 0 e 60 no formato "1-2-3-4-5-6" 
Gerar sequências de 6 números aleatorios e comparar com os números informados
Quando baterem os números(a geração não precisa respeitar a ordem dos números) o programa deve ser encerrado informando:
* Número de tentativas para bater a sequencia informada
* Tempo de execução em millisegundos
* Numero de sequencias geradas por ms (tentativas/tempo)
* Sequencia aleatória que foi gerada

Objetivo é gerar a maior quantidade de sequencias aleatórias por millisegundos
(o melhor ganha uma pizza)
ex.: ./programa "0-15-16-20-21-60"

## Método

Foi utilizado o algoritmo **xorshift+** como gerador de números aleatórios e com 5 workers (goroutines) para adivinhação. A validação foi feita com base no checksum dos números, comparando os arrays apenas com checksums iguais.

## build

```bash
go build
```

## Exemplo de uso

```bash
./sena-challenge 1-2-3-4-5-6
```
