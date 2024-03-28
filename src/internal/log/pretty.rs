use crate::internal::log::json;

pub fn print(string: String) {
    if json::JSON_LOGGING {return}

    println!("{}", string)
}