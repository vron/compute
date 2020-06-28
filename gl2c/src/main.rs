/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */
use std::fs::File;
use std::io::Write;
use gl2c::translate;

fn main() {

    // assume the first arguemnt is the file to translate, the second where we
    // should write the cpp file to


    let mut a = std::env::args();
    if a.len() < 3 {
        panic!("must have input, output as args")
    }
    let _ = a.next();
    let inf = a.next().unwrap();
    let name =  a.next().unwrap();
    let mut file = File::create(format!("{}", name)).expect("could not open file");
    let data = translate(inf);
    let _ = write!(file, "{}", data);
}
