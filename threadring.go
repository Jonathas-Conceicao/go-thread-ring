/**
 * Thread Ring Problem from 'The Computer Language Benchmarks Game'.
 * http://benchmarksgame.alioth.debian.org/u64q/threadring-description.html#threadring
 *
 * by: Jonathas Conceição
 * email: jadoliveira@inf.ufpel.edu.br
 */

package main

import "fmt"
import "os"
import "strconv"
import "runtime"
import "sync"

type Thread struct { // Struct que representa uma thread em execução.
  id   int           // Id da Thread.
  Op   int           // O tolken que é passado de thread para thread que também é o número de operações restantes.
  mtx sync.Mutex     // Toda thread carrega um mutex que é adiquirido e mantido assim até que ela receba o token.
  next *Thread       // Referencia para a próxima thread.
}

/**
 * pass() será executado num ponto onde o mutex da thread estiver desbloqueado.
 * Se o token tiver o valor zero, então o ID da thread é enviado para o 'channel done'.
 * Se o token tiver outro valor, o valor é decrementado, enviado para a próxima thread e então a desbloqueia.
 */
func (t *Thread) pass() {
  if t.Op == 0 {
    done <- t.id // Envia o ID da thread para o channel.
  } else {
    t.next.Op = t.Op -1
    t.next.mtx.Unlock() // Librea o lock da próxima thread para que ela possa processar o token.
  }
}

/**
 * run() é o loop principal de cada thread, quando ela obtém o lock com sucesso
 * ela executa o método run() para passar o token adiante.
 */
func (t *Thread) run() {
  for {
    t.mtx.Lock() // Bloqueia até que o lock possa ser obtido.
    t.pass()
    runtime.Gosched() // Libera o processador permitindo que outras goroutines sejam executadas.
  }
}

/**
 * Inicia uma thread e adquire o lock do seu mutex e dispara a goroutine.
 */
func (t *Thread) Init(id int, n *Thread) {
  t.id   = id
  t.next = n
  t.mtx.Lock()
  go t.run() // Dispara t.run() em uma nova goroutine (lightweight thread).
}

/**
 * Passa o token inicial (número de operações) para uma thread
 * e librea seu lock para que a passagem de token começe.
 */
func (t *Thread) Start(ops int) {
  t.Op = ops
  t.mtx.Unlock()
}

const nThreads = 503

var done chan int = make(chan int) // cria um channel usado para anunciar o fim do programa.

func main() {
  ops, _ := strconv.Atoi(os.Args[1])
  runtime.GOMAXPROCS(1) // Limita o programa ao uso de uma CPU
  var threads [nThreads]Thread
  for i := range(threads) {
    threads[i].Init(i+1, &threads[(i+1) % nThreads])
  }
  threads[0].Start(ops)
  fmt.Println(<-done) // espera até que um dado esteja disponivel no channel, então o retira, printa e termina a execução.
}
